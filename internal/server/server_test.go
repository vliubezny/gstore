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
