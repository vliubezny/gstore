package server

import (
	"context"
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
			desc:     "invalid: missing email",
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
