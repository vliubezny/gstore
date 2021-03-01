package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

//go:generate mockgen -destination=./service_mock.go -package=auth -source=service.go

const (
	typeAccess  = "access"
	typeRefresh = "refresh"

	issuer = "gstore.auth"
)

var (
	// ErrEmailIsTaken states that email address is taken.
	ErrEmailIsTaken = errors.New("email is taken")

	// ErrInvalidCredentials states that email or password is invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidToken states that token is invalid.
	ErrInvalidToken = errors.New("invalid token")
)

// AccessTokenClaims specifies the claims for access token.
type AccessTokenClaims struct {
	TokenType string `json:"type,omitempty"`
	UserID    int64  `json:"userId,omitempty"`
	IsAdmin   bool   `json:"admin,omitempty"`
	jwt.StandardClaims
}

// RefreshTokenClaims specifies the claims for access token.
type RefreshTokenClaims struct {
	TokenType string `json:"type,omitempty"`
	UserID    int64  `json:"userId,omitempty"`
	jwt.StandardClaims
}

// TokenPair groups access and refresh tokens.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// AccessTokenValidator parses and validates access token.
type AccessTokenValidator func(token string) (AccessTokenClaims, error)

// Service provides methods for user authentication.
type Service interface {
	Register(ctx context.Context, user model.User, password string) (model.User, error)
	Login(ctx context.Context, email, password string) (TokenPair, error)
	Refresh(ctx context.Context, refreshToken string) (TokenPair, error)
	Revoke(ctx context.Context, refreshToken string) error
	ValidateAccessToken(token string) (AccessTokenClaims, error)
}

type authService struct {
	s       storage.UserStorage
	signKey []byte
}

// New creates instance of auth service.
func New(s storage.UserStorage, signKey string) Service {
	return &authService{
		s:       s,
		signKey: []byte(signKey),
	}
}

func (s *authService) Register(ctx context.Context, user model.User, password string) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hash)

	user, err = s.s.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrEmailIsTaken) {
			return model.User{}, ErrEmailIsTaken
		}

		return model.User{}, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (TokenPair, error) {
	u, err := s.s.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return TokenPair{}, ErrInvalidCredentials
		}
		return TokenPair{}, fmt.Errorf("failed to get user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return TokenPair{}, ErrInvalidCredentials
	}

	at, err := s.signToken(newAccessClaims(u))
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	rc := newRefreshClaims(u)

	if err = s.s.SaveToken(ctx, rc.Id, u.ID, time.Unix(rc.ExpiresAt, 0)); err != nil {
		return TokenPair{}, fmt.Errorf("failed to save refresh token: %w", err)
	}

	rt, err := s.signToken(rc)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return TokenPair{AccessToken: at, RefreshToken: rt}, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (TokenPair, error) {
	claims, err := validateRefreshToken(refreshToken, s.signKey)
	if err != nil {
		return TokenPair{}, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	u, err := s.s.GetUserByID(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return TokenPair{}, fmt.Errorf("%w: missing user", ErrInvalidToken)
		}
		return TokenPair{}, fmt.Errorf("failed to get user: %w", err)
	}

	if err = s.s.DeleteToken(ctx, claims.Id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return TokenPair{}, fmt.Errorf("%w: token has been used", ErrInvalidToken)
		}
		return TokenPair{}, fmt.Errorf("failed to delete token: %w", err)
	}

	at, err := s.signToken(newAccessClaims(u))
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign access token: %w", err)
	}

	rc := newRefreshClaims(u)

	if err = s.s.SaveToken(ctx, rc.Id, u.ID, time.Unix(rc.ExpiresAt, 0)); err != nil {
		return TokenPair{}, fmt.Errorf("failed to save refresh token: %w", err)
	}

	rt, err := s.signToken(rc)
	if err != nil {
		return TokenPair{}, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return TokenPair{AccessToken: at, RefreshToken: rt}, nil
}

func (s *authService) Revoke(ctx context.Context, refreshToken string) error {
	claims, err := validateRefreshToken(refreshToken, s.signKey)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if err = s.s.DeleteToken(ctx, claims.Id); err != nil && !errors.Is(err, storage.ErrNotFound) {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (s *authService) signToken(claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.signKey)
}

func newAccessClaims(user model.User) AccessTokenClaims {
	return AccessTokenClaims{
		TokenType: typeAccess,
		UserID:    user.ID,
		IsAdmin:   user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
}

func (s *authService) ValidateAccessToken(token string) (AccessTokenClaims, error) {
	at, err := jwt.ParseWithClaims(token, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("token must be signed with HS256 alg")
		}
		return s.signKey, nil
	})

	if err != nil {
		return AccessTokenClaims{}, fmt.Errorf("unable to parse claims: %w", err)
	}

	claims, ok := at.Claims.(*AccessTokenClaims)
	if !ok || !at.Valid {
		return AccessTokenClaims{}, errors.New("invalid token")
	}
	if claims.TokenType != typeAccess {
		return AccessTokenClaims{}, fmt.Errorf("invalid access token: type %s", claims.TokenType)
	}
	return *claims, nil
}

func newRefreshClaims(user model.User) RefreshTokenClaims {
	return RefreshTokenClaims{
		TokenType: typeRefresh,
		UserID:    user.ID,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
		},
	}
}

func validateRefreshToken(token string, signKey []byte) (RefreshTokenClaims, error) {
	at, err := jwt.ParseWithClaims(token, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("token must be signed with HS256 alg")
		}
		return signKey, nil
	})

	if err != nil {
		return RefreshTokenClaims{}, fmt.Errorf("unable to parse claims: %w", err)
	}

	claims, ok := at.Claims.(*RefreshTokenClaims)
	if !ok || !at.Valid {
		return RefreshTokenClaims{}, errors.New("invalid token")
	}
	if claims.TokenType != typeRefresh {
		return RefreshTokenClaims{}, fmt.Errorf("invalid refresh token: type %s", claims.TokenType)
	}
	return *claims, nil
}
