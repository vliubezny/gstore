package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

var (
	// ErrEmailIsTaken states that email address is taken.
	ErrEmailIsTaken = errors.New("email is taken")

	// ErrInvalidCredentials states that email or password is invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

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

	return createToken(u, s.signKey)
}

func createToken(user model.User, signKey string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = user.ID
	claims["admin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return at.SignedString([]byte(signKey))
}
