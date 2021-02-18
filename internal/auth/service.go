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
)

var (
	// ErrEmailIsTaken states that email address is taken.
	ErrEmailIsTaken = errors.New("email is taken")

	// ErrInvalidCredentials states that email or password is invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
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

// Service provides methods for user authentication.
type Service interface {
	Register(ctx context.Context, user model.User, password string) (model.User, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	s       storage.UserStorage
	signKey string
}

// New creates instance of auth service.
func New(s storage.UserStorage, signKey string) Service {
	return &authService{
		s:       s,
		signKey: signKey,
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

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.s.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return createAccessToken(u, s.signKey)
}

func createAccessToken(user model.User, signKey string) (string, error) {
	claims := AccessTokenClaims{
		TokenType: typeAccess,
		UserID:    user.ID,
		IsAdmin:   user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			Issuer:    "gstore.auth",
			ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		},
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(signKey))
}

func validateAccessToken(token, signKey string) (AccessTokenClaims, error) {
	at, err := jwt.ParseWithClaims(token, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("token must be signed with HS256 alg")
		}
		return []byte(signKey), nil
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
