package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gstore/internal/auth"
	"github.com/vliubezny/gstore/internal/model"
)

func Test_registerHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		user     model.User
		password string
		err      error
		input    string
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			user:     model.User{Email: "admin@test.com"},
			password: "testP@ss",
			err:      nil,
			input:    `{"email":"admin@test.com", "password":"testP@ss"}`,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "email":"admin@test.com", "isAdmin":false}`,
		},
		{
			desc:     "invalid: missing email",
			user:     model.User{},
			password: "testP@ss",
			err:      errSkip,
			input:    `{"password":"testP@ss"}`,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error": "email is a required field"}`,
		},
		{
			desc:     "invalid: taken email",
			user:     model.User{Email: "admin@test.com"},
			password: "testP@ss",
			err:      auth.ErrEmailIsTaken,
			input:    `{"email":"admin@test.com", "password":"testP@ss"}`,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error": "email address has been already taken"}`,
		},
		{
			desc:     "internal error",
			user:     model.User{Email: "admin@test.com"},
			password: "testP@ss",
			err:      errTest,
			input:    `{"email":"admin@test.com", "password":"testP@ss"}`,
			rcode:    http.StatusInternalServerError,
			rdata:    `{"error": "internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := auth.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().Register(gomock.Any(), tC.user, tC.password).
					DoAndReturn(func(_ context.Context, u model.User, _ string) (model.User, error) {
						u.ID = 1
						return u, tC.err
					})
			}

			router := setupTestRouterWithAuth(nil, svc)
			rec, r := newTestParameters(http.MethodPost, "/v1/register", tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_loginHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		email    string
		password string
		tokens   auth.TokenPair
		err      error
		input    string
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			email:    "admin@test.com",
			password: "testP@ss",
			tokens:   auth.TokenPair{AccessToken: "testAccess", RefreshToken: "testRefresh"},
			err:      nil,
			input:    `{"email":"admin@test.com", "password":"testP@ss"}`,
			rcode:    http.StatusOK,
			rdata:    `{"accessToken":"testAccess", "refreshToken":"testRefresh"}`,
		},
		{
			desc:     "invalid credentials",
			email:    "admin@test.com",
			password: "testP@ss",
			tokens:   auth.TokenPair{},
			err:      auth.ErrInvalidCredentials,
			input:    `{"email":"admin@test.com", "password":"testP@ss"}`,
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"invalid username or password"}`,
		},
		{
			desc:     "internal error",
			email:    "admin@test.com",
			password: "testP@ss",
			tokens:   auth.TokenPair{},
			err:      assert.AnError,
			input:    `{"email":"admin@test.com", "password":"testP@ss"}`,
			rcode:    http.StatusInternalServerError,
			rdata:    `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := auth.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().Login(gomock.Any(), tC.email, tC.password).Return(tC.tokens, tC.err)
			}

			router := setupTestRouterWithAuth(nil, svc)
			rec, r := newTestParameters(http.MethodPost, "/v1/login", tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_refreshHandler(t *testing.T) {
	testCases := []struct {
		desc   string
		token  string
		tokens auth.TokenPair
		err    error
		rcode  int
		rdata  string
	}{
		{
			desc:   "success",
			token:  "testToken",
			tokens: auth.TokenPair{AccessToken: "testAccess", RefreshToken: "testRefresh"},
			err:    nil,
			rcode:  http.StatusOK,
			rdata:  `{"accessToken":"testAccess", "refreshToken":"testRefresh"}`,
		},
		{
			desc:   "missing token - Unauthorized",
			token:  "",
			tokens: auth.TokenPair{AccessToken: "testAccess", RefreshToken: "testRefresh"},
			err:    errSkip,
			rcode:  http.StatusUnauthorized,
			rdata:  `{"error":"missing token"}`,
		},
		{
			desc:   "invalid token - Unauthorized",
			token:  "testtoken",
			tokens: auth.TokenPair{AccessToken: "testAccess", RefreshToken: "testRefresh"},
			err:    auth.ErrInvalidToken,
			rcode:  http.StatusUnauthorized,
			rdata:  `{"error":"invalid refresh token"}`,
		},
		{
			desc:   "error",
			token:  "testtoken",
			tokens: auth.TokenPair{AccessToken: "testAccess", RefreshToken: "testRefresh"},
			err:    assert.AnError,
			rcode:  http.StatusInternalServerError,
			rdata:  `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := auth.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().Refresh(gomock.Any(), tC.token).Return(tC.tokens, tC.err)
			}

			router := setupTestRouterWithAuth(nil, svc)
			rec, r := newTestParameters(http.MethodPost, "/v1/refresh", "")
			r.Header.Del("Authorization")

			if tC.token != "" {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tC.token))
			}

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_revokeHandler(t *testing.T) {
	testCases := []struct {
		desc  string
		token string
		err   error
		rcode int
		rdata string
	}{
		{
			desc:  "success",
			token: "testToken",
			err:   nil,
			rcode: http.StatusNoContent,
			rdata: "",
		},
		{
			desc:  "missing token - Unauthorized",
			token: "",
			err:   errSkip,
			rcode: http.StatusUnauthorized,
			rdata: `{"error":"missing token"}`,
		},
		{
			desc:  "invalid token - Unauthorized",
			token: "testtoken",
			err:   auth.ErrInvalidToken,
			rcode: http.StatusUnauthorized,
			rdata: `{"error":"invalid refresh token"}`,
		},
		{
			desc:  "error",
			token: "testtoken",
			err:   assert.AnError,
			rcode: http.StatusInternalServerError,
			rdata: `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := auth.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().Revoke(gomock.Any(), tC.token).Return(tC.err)
			}

			router := setupTestRouterWithAuth(nil, svc)
			rec, r := newTestParameters(http.MethodPost, "/v1/revoke", "")
			r.Header.Del("Authorization")

			if tC.token != "" {
				r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tC.token))
			}

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			if tC.rdata == "" {
				assert.Empty(t, string(body))
			} else {
				assert.JSONEq(t, tC.rdata, string(body))
			}
		})
	}
}

func Test_updateUserPermissionsHandler(t *testing.T) {
	testCases := []struct {
		desc  string
		id    string
		user  model.User
		err   error
		input string
		rcode int
		rdata string
	}{
		{
			desc:  "success",
			id:    "1",
			user:  model.User{ID: 1, IsAdmin: true},
			err:   nil,
			input: `{"isAdmin":true}`,
			rcode: http.StatusNoContent,
			rdata: ``,
		},
		{
			desc:  "invalid id",
			id:    "test",
			user:  model.User{ID: 1, IsAdmin: true},
			err:   errSkip,
			input: `{"isAdmin":true}`,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"invalid user ID"}`,
		},
		{
			desc:  "invalid payload",
			id:    "1",
			user:  model.User{ID: 1, IsAdmin: true},
			err:   errSkip,
			input: `{"admin":true}`,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"isAdmin is a required field"}`,
		},
		{
			desc:  "user not found",
			id:    "1",
			user:  model.User{ID: 1, IsAdmin: true},
			err:   auth.ErrNotFound,
			input: `{"isAdmin":true}`,
			rcode: http.StatusNotFound,
			rdata: `{"error":"user not found"}`,
		},
		{
			desc:  "internal error",
			id:    "1",
			user:  model.User{ID: 1, IsAdmin: true},
			err:   assert.AnError,
			input: `{"isAdmin":true}`,
			rcode: http.StatusInternalServerError,
			rdata: `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := auth.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().UpdateUserPermissions(gomock.Any(), tC.user).Return(tC.err)
			}

			router := setupTestRouterWithAuth(nil, svc)
			rec, r := newTestParameters(http.MethodPut, fmt.Sprintf("/v1/users/%s/permissions", tC.id), tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			if tC.rdata == "" {
				assert.Empty(t, string(body))
			} else {
				assert.JSONEq(t, tC.rdata, string(body))
			}
		})
	}
}
