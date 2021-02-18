package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	signKey  = "secret"
	testPass = "test123"
	testHash = "$2a$10$Ej1ANHun0jp1O5ozBhTbGODKprti6Z2FheUyHdyuvcJ6/feFo9s/K"
)

var ctx = context.Background()

func TestService_Register(t *testing.T) {
	testCases := []struct {
		desc     string
		rErr     error
		user     model.User
		password string
		err      error
	}{
		{
			desc:     "success",
			rErr:     nil,
			user:     model.User{Email: "admin@test.com", IsAdmin: true},
			password: testPass,
			err:      nil,
		},
		{
			desc:     "ErrEmailIsTaken",
			rErr:     storage.ErrEmailIsTaken,
			user:     model.User{Email: "admin@test.com", IsAdmin: true},
			password: testPass,
			err:      ErrEmailIsTaken,
		},
		{
			desc:     "unexpected error",
			rErr:     assert.AnError,
			user:     model.User{Email: "admin@test.com", IsAdmin: true},
			password: testPass,
			err:      assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockUserStorage(ctrl)

			st.EXPECT().CreateUser(ctx, gomock.AssignableToTypeOf(model.User{})).
				DoAndReturn(func(_ context.Context, u model.User) (model.User, error) {
					u.ID = 1
					return u, tC.rErr
				})

			s := New(st, signKey)

			u, err := s.Register(ctx, tC.user, tC.password)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			if err == nil {
				assert.Equal(t, int64(1), u.ID)
				assert.Equal(t, tC.user.Email, u.Email)
				assert.Equal(t, tC.user.IsAdmin, u.IsAdmin)
				assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(tC.password)), "incorrect password hash")
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	testCases := []struct {
		desc     string
		rUser    model.User
		rErr     error
		email    string
		password string
		err      error
	}{
		{
			desc:     "success",
			rUser:    model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rErr:     nil,
			email:    "admin@test.com",
			password: testPass,
			err:      nil,
		},
		{
			desc:     "invelid email",
			rUser:    model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rErr:     storage.ErrNotFound,
			email:    "admin@test.com",
			password: testPass,
			err:      ErrInvalidCredentials,
		},
		{
			desc:     "invelid password",
			rUser:    model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rErr:     nil,
			email:    "admin@test.com",
			password: "invalid",
			err:      ErrInvalidCredentials,
		},
		{
			desc:     "error",
			rUser:    model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rErr:     assert.AnError,
			email:    "admin@test.com",
			password: testPass,
			err:      assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockUserStorage(ctrl)

			st.EXPECT().GetUserByEmail(ctx, tC.email).Return(tC.rUser, tC.rErr)

			s := New(st, signKey)

			token, err := s.Login(ctx, tC.email, tC.password)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			if err == nil {
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestService_createAccessToken(t *testing.T) {
	u := model.User{
		ID:      1,
		Email:   "admin@test.com",
		IsAdmin: true,
	}

	at, err := createAccessToken(u, signKey)

	require.NoError(t, err)
	require.NotEmpty(t, at)

	claims, err := validateAccessToken(at, signKey)

	require.NoError(t, err)
	assert.Equal(t, u.ID, claims.UserID)
	assert.Equal(t, u.IsAdmin, claims.IsAdmin)
}
