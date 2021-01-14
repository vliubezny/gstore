package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vliubezny/gstore/internal/auth"
)

func Test_setContentTypeMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", nil)

	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":42}`))
	})

	setContentTypeMiddleware(contentTypeJSON)(h).ServeHTTP(rec, r)
	assert.Equal(t, "application/json", rec.Result().Header.Get("Content-Type"))
}

func Test_loggerMiddleware(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	hook := test.NewGlobal()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
	req.Header.Set("User-Agent", "curl")
	req.Header.Set("X-Forwarded-For", "210.172.60.240")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})

	loggerMiddleware(h).ServeHTTP(rec, req)

	log := hook.LastEntry()
	require.NotNil(t, log)
	assert.Equal(t, logrus.DebugLevel, log.Level)
	assert.Equal(t, "POST /v1/test", log.Message, "Incorrect request entry")
	assert.Equal(t, "curl", log.Data["agent"], "Incorrect user agent")
	assert.Equal(t, "210.172.60.240", log.Data["ip"], "Incorrect user IP")
}

func Test_recoveryMiddleware(t *testing.T) {
	logger, hook := test.NewNullLogger()
	ctx := context.WithValue(context.Background(), loggerKey{}, logger)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("", "/", nil).WithContext(ctx)

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	recoveryMiddleware(h).ServeHTTP(rec, req)

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	assert.Equal(t, `{"error":"internal error"}`, string(body))

	log := hook.LastEntry()
	require.NotNil(t, log)
	assert.Equal(t, logrus.ErrorLevel, log.Level)
	assert.Contains(t, log.Message, "test panic", "Missing panic message")
	assert.Contains(t, log.Message, file, "Missing stacktrace")
}

func Test_basicAuthMiddleware(t *testing.T) {
	username := "admin"
	password := "pass"
	testCases := []struct {
		desc     string
		method   string
		withAuth bool
		valid    bool
		err      error
		rcode    int
		rdata    string
	}{
		{
			desc:     "GET: allow anonymous",
			method:   http.MethodGet,
			withAuth: false,
			valid:    false,
			err:      errSkip,
			rcode:    http.StatusOK,
			rdata:    `{"result":"OK"}`,
		},
		{
			desc:     "POST: block anonymous",
			method:   http.MethodPost,
			withAuth: false,
			valid:    false,
			err:      errSkip,
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
		{
			desc:     "PUT: block anonymous",
			method:   http.MethodPut,
			withAuth: false,
			valid:    false,
			err:      errSkip,
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
		{
			desc:     "DELETE: block anonymous",
			method:   http.MethodDelete,
			withAuth: false,
			valid:    false,
			err:      errSkip,
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
		{
			desc:     "POST: block invalid credentials",
			method:   http.MethodPost,
			withAuth: true,
			valid:    false,
			err:      nil,
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
		{
			desc:     "POST: allow valid credentials",
			method:   http.MethodPost,
			withAuth: true,
			valid:    true,
			err:      nil,
			rcode:    http.StatusOK,
			rdata:    `{"result":"OK"}`,
		},
		{
			desc:     "POST: auth internal error",
			method:   http.MethodPost,
			withAuth: true,
			valid:    false,
			err:      errTest,
			rcode:    http.StatusInternalServerError,
			rdata:    `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			logger, _ := test.NewNullLogger()
			ctx := context.WithValue(context.Background(), loggerKey{}, logger)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tC.method, "/", nil).WithContext(ctx)

			if tC.withAuth {
				req.SetBasicAuth(username, password)
			}

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"result":"OK"}`))
			})

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			a := auth.NewMockAuthenticator(ctrl)
			if tC.err != errSkip {
				a.EXPECT().Authenticate(username, password).Return(tC.valid, tC.err)
			}

			basicAuthMiddleware(a)(h).ServeHTTP(rec, req)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}
