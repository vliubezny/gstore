package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
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

func mustCreateAccessToken(u model.User) string {
	c := newAccessClaims(u)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(signKey))
	if err != nil {
		panic(err)
	}
	return token
}

func mustCreateRefreshToken(u model.User) string {
	c := newRefreshClaims(u)
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString([]byte(signKey))
	if err != nil {
		panic(err)
	}
	return token
}

func TestService_Refresh(t *testing.T) {
	user := model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true}
	testCases := []struct {
		desc            string
		rUser           model.User
		rUserErr        error
		rDeleteTokenErr error
		rSaveTokenErr   error
		token           string
		err             error
	}{
		{
			desc:            "success",
			rUser:           user,
			rUserErr:        nil,
			rDeleteTokenErr: nil,
			rSaveTokenErr:   nil,
			token:           mustCreateRefreshToken(user),
			err:             nil,
		},
		{
			desc:            "malformed token - ErrInvalidToken",
			rUser:           model.User{},
			rUserErr:        errSkip,
			rDeleteTokenErr: errSkip,
			rSaveTokenErr:   errSkip,
			token:           "test",
			err:             ErrInvalidToken,
		},
		{
			desc:            "missing user - ErrInvalidToken",
			rUser:           user,
			rUserErr:        storage.ErrNotFound,
			rDeleteTokenErr: errSkip,
			rSaveTokenErr:   errSkip,
			token:           mustCreateRefreshToken(user),
			err:             ErrInvalidToken,
		},
		{
			desc:            "get user - error",
			rUser:           user,
			rUserErr:        assert.AnError,
			rDeleteTokenErr: errSkip,
			rSaveTokenErr:   errSkip,
			token:           mustCreateRefreshToken(user),
			err:             assert.AnError,
		},
		{
			desc:            "token used - ErrInvalidToken",
			rUser:           user,
			rUserErr:        nil,
			rDeleteTokenErr: storage.ErrNotFound,
			rSaveTokenErr:   errSkip,
			token:           mustCreateRefreshToken(user),
			err:             ErrInvalidToken,
		},
		{
			desc:            "delete token - error",
			rUser:           user,
			rUserErr:        nil,
			rDeleteTokenErr: assert.AnError,
			rSaveTokenErr:   errSkip,
			token:           mustCreateRefreshToken(user),
			err:             assert.AnError,
		},
		{
			desc:            "save token - error",
			rUser:           user,
			rUserErr:        nil,
			rDeleteTokenErr: nil,
			rSaveTokenErr:   assert.AnError,
			token:           mustCreateRefreshToken(user),
			err:             assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockUserStorage(ctrl)
			tx := storage.NewMockUserStorage(ctrl)

			if tC.rUserErr != errSkip {
				st.EXPECT().GetUserByID(ctx, tC.rUser.ID).Return(tC.rUser, tC.rUserErr)
			}

			if tC.rDeleteTokenErr != errSkip {
				st.EXPECT().InTx(ctx, gomock.Any()).DoAndReturn(
					func(_ context.Context, action func(s storage.UserStorage) error) error {
						return action(tx)
					})
				tx.EXPECT().DeleteToken(ctx, gomock.AssignableToTypeOf("")).Return(tC.rDeleteTokenErr)
			}

			if tC.rSaveTokenErr != errSkip {
				tx.EXPECT().SaveToken(ctx, gomock.AssignableToTypeOf(""), tC.rUser.ID, gomock.AssignableToTypeOf(time.Time{})).
					Return(tC.rSaveTokenErr)
			}

			s := New(st, signKey)

			pair, err := s.Refresh(ctx, tC.token)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			if err == nil {
				assert.NotEmpty(t, pair.AccessToken)
				assert.NotEmpty(t, pair.RefreshToken)
			}
		})
	}
}

func TestService_Revoke(t *testing.T) {
	user := model.User{ID: 1, Email: "admin@test.com", PasswordHash: testHash, IsAdmin: true}
	testCases := []struct {
		desc  string
		rErr  error
		token string
		err   error
	}{
		{
			desc:  "success",
			rErr:  nil,
			token: mustCreateRefreshToken(user),
			err:   nil,
		},
		{
			desc:  "token used - success",
			rErr:  storage.ErrNotFound,
			token: mustCreateRefreshToken(user),
			err:   nil,
		},
		{
			desc:  "malformed token - ErrInvalidToken",
			rErr:  errSkip,
			token: "test",
			err:   ErrInvalidToken,
		},
		{
			desc:  "delete token - error",
			rErr:  assert.AnError,
			token: mustCreateRefreshToken(user),
			err:   assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockUserStorage(ctrl)

			if tC.rErr != errSkip {
				st.EXPECT().DeleteToken(ctx, gomock.AssignableToTypeOf("")).Return(tC.rErr)
			}

			s := New(st, signKey)

			err := s.Revoke(ctx, tC.token)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_newAccessClaims(t *testing.T) {
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
	assert.InDelta(t, time.Now().Add(10*time.Minute).Unix(), ac.ExpiresAt, 1)
}

func TestService_newRefreshClaims(t *testing.T) {
	u := model.User{
		ID:      1,
		Email:   "admin@test.com",
		IsAdmin: true,
	}

	ac := newRefreshClaims(u)

	assert.Equal(t, typeRefresh, ac.TokenType)
	assert.Equal(t, u.ID, ac.UserID)
	assert.NotEmpty(t, ac.Id)
	assert.InDelta(t, time.Now().Add(30*24*time.Hour).Unix(), ac.ExpiresAt, 1)
}

func TestService_validateRefreshToken(t *testing.T) {
	s := &authService{
		signKey: []byte(signKey),
	}

	u := model.User{
		ID:      1,
		Email:   "admin@test.com",
		IsAdmin: true,
	}

	ac := newRefreshClaims(u)
	token, err := s.signToken(ac)
	require.NoError(t, err)

	ac, err = validateRefreshToken(token, []byte(signKey))
	require.NoError(t, err)

	assert.Equal(t, typeRefresh, ac.TokenType)
	assert.Equal(t, u.ID, ac.UserID)
	assert.NotEmpty(t, ac.Id)
	assert.InDelta(t, time.Now().Add(30*24*time.Hour).Unix(), ac.ExpiresAt, 1)
}

func TestService_ValidateAccessToken(t *testing.T) {
	s := New(nil, signKey)

	u := model.User{
		ID:      1,
		Email:   "admin@test.com",
		IsAdmin: true,
	}

	token := mustCreateAccessToken(u)

	claims, err := s.ValidateAccessToken(token)
	require.NoError(t, err)

	assert.Equal(t, typeAccess, claims.TokenType)
	assert.Equal(t, u.ID, claims.UserID)
	assert.NotEmpty(t, claims.Id)
}

func TestService_UpdateUserPermissions(t *testing.T) {
	testCases := []struct {
		desc string
		user model.User
		rErr error
		err  error
	}{
		{
			desc: "success",
			user: model.User{ID: 1, IsAdmin: true},
			rErr: nil,
			err:  nil,
		},
		{
			desc: "ErrNotFound",
			rErr: storage.ErrNotFound,
			user: model.User{ID: 1, IsAdmin: true},
			err:  ErrNotFound,
		},
		{
			desc: "unexpected error",
			rErr: assert.AnError,
			user: model.User{ID: 1, IsAdmin: true},
			err:  assert.AnError,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockUserStorage(ctrl)

			st.EXPECT().UpdateUserPermissions(ctx, tC.user).Return(tC.rErr)

			s := New(st, signKey)

			err := s.UpdateUserPermissions(ctx, tC.user)

			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}
