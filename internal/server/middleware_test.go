package server

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
