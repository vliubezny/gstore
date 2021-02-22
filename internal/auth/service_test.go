package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	signKey  = "secret"
	testPass = "test123"
	testHash = "$2a$10$Ej1ANHun0jp1O5ozBhTbGODKprti6Z2FheUyHdyuvcJ6/feFo9s/K"
)

var (
	ctx     = context.Background()
	errSkip = errors.New("skip")
)

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
		desc      string
		rUser     model.User
		rUserErr  error
		rTokenErr error
		email     string
		password  string
		err       error
	}{
		{
			desc:      "success",
			rUser:     model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rUserErr:  nil,
			rTokenErr: nil,
			email:     "admin@test.com",
			password:  testPass,
			err:       nil,
		},
		{
			desc:      "invalid email",
			rUser:     model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rUserErr:  storage.ErrNotFound,
			rTokenErr: errSkip,
			email:     "admin@test.com",
			password:  testPass,
			err:       ErrInvalidCredentials,
		},
		{
			desc:      "invalid password",
			rUser:     model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rUserErr:  nil,
			rTokenErr: errSkip,
			email:     "admin@test.com",
			password:  "invalid",
			err:       ErrInvalidCredentials,
		},
		{
			desc:      "error getting user",
			rUser:     model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rUserErr:  assert.AnError,
			rTokenErr: errSkip,
			email:     "admin@test.com",
			password:  testPass,
			err:       assert.AnError,
		},
		{
			desc:      "error saving token",
			rUser:     model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true},
			rUserErr:  nil,
			rTokenErr: assert.AnError,
			email:     "admin@test.com",
			password:  testPass,
			err:       assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockUserStorage(ctrl)

			st.EXPECT().GetUserByEmail(ctx, tC.email).Return(tC.rUser, tC.rUserErr)

			if tC.rTokenErr != errSkip {
				st.EXPECT().SaveToken(ctx, gomock.AssignableToTypeOf(""), tC.rUser.ID, gomock.AssignableToTypeOf(time.Time{})).
					Return(tC.rTokenErr)
			}

			s := New(st, signKey)

			pair, err := s.Login(ctx, tC.email, tC.password)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			if err == nil {
				assert.NotEmpty(t, pair.AccessToken)
				assert.NotEmpty(t, pair.RefreshToken)
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

	ac := newAccessClaims(u)

	assert.Equal(t, typeAccess, ac.TokenType)
	assert.Equal(t, u.ID, ac.UserID)
	assert.Equal(t, u.IsAdmin, ac.IsAdmin)
	assert.NotEmpty(t, ac.Id)
}
