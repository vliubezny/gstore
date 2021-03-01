package server

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vliubezny/gstore/internal/auth"
	"github.com/vliubezny/gstore/internal/service"
)

const (
	testUsername = "admin"
	testPassword = "pass123"
)

func setupTestRouter(s service.Service) http.Handler {
	return setupTestRouterWithAuth(s, nil)
}

func setupTestRouterWithAuth(s service.Service, a auth.Service) http.Handler {
	r := chi.NewRouter()
	SetupRouter(s, a, r, func(_ string) (auth.AccessTokenClaims, error) {
		return auth.AccessTokenClaims{UserID: 1, IsAdmin: true}, nil
	})
	return r
}

func newTestParameters(method, uri, body string) (*httptest.ResponseRecorder, *http.Request) {
	test.NewGlobal()
	rec := httptest.NewRecorder()

	var payload io.Reader
	if body != "" {
		payload = strings.NewReader(body)
	}
	r := httptest.NewRequest(method, uri, payload)
	if body != "" {
		r.Header.Set(headerContentType, contentTypeJSON)
	}

	r.Header.Set("Authorization", "Bearer testtoken")

	return rec, r
}

func Test_getLogger(t *testing.T) {
	l := logrus.New()
	ctx := context.WithValue(context.Background(), loggerKey{}, l)
	r := httptest.NewRequest("", "/", nil).WithContext(ctx)

	logger := getLogger(r)

	assert.Exactly(t, l, logger)
}

func Test_writeError(t *testing.T) {
	logger, hook := test.NewNullLogger()
	rec := httptest.NewRecorder()

	writeError(logger, rec, http.StatusBadRequest, "test error")

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
	assert.Equal(t, `{"error":"test error"}`, string(body))

	log := hook.LastEntry()
	require.NotNil(t, log)
	assert.Equal(t, logrus.ErrorLevel, log.Level)
	assert.Contains(t, log.Message, "test error", "Missing error message")
}

func Test_writeInternalError(t *testing.T) {
	logger, hook := test.NewNullLogger()
	rec := httptest.NewRecorder()

	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)

	writeInternalError(logger, rec, "test error")

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
	assert.Equal(t, `{"error":"internal error"}`, string(body))

	log := hook.LastEntry()
	require.NotNil(t, log)
	assert.Equal(t, logrus.ErrorLevel, log.Level)
	assert.Contains(t, log.Message, "test error", "Missing error message")
	assert.Contains(t, log.Message, file, "Missing stacktrace")
}

func Test_writeOK(t *testing.T) {
	logger, _ := test.NewNullLogger()
	rec := httptest.NewRecorder()

	writeOK(logger, rec, struct {
		Msg string
	}{
		Msg: "test",
	})

	body, _ := ioutil.ReadAll(rec.Result().Body)

	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, `{"Msg":"test"}`, string(body))
}

func Test_extractBearer(t *testing.T) {
	testCases := []struct {
		desc  string
		auth  string
		token string
	}{
		{
			desc:  "success",
			auth:  "Bearer testtoken",
			token: "testtoken",
		},
		{
			desc:  "missing header",
			auth:  "",
			token: "",
		},
		{
			desc:  "incorrect header",
			auth:  "Basic 1223232323",
			token: "",
		},
		{
			desc:  "empty bearer",
			auth:  "Bearer ",
			token: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r := httptest.NewRequest("", "/", nil)
			if tC.auth != "" {
				r.Header.Set("Authorization", tC.auth)
			}

			token := extractBearer(r)

			assert.Equal(t, tC.token, token)
		})
	}
}
