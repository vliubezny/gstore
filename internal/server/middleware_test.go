package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
	b := bytes.NewBufferString("")
	logrus.SetOutput(b)
	logrus.SetLevel(logrus.DebugLevel)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
	req.Header.Set("User-Agent", "curl")
	req.Header.Set("X-Forwarded-For", "210.172.60.240")

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(nil)
	})

	loggerMiddleware(h).ServeHTTP(rec, req)

	log := b.String()
	assert.Contains(t, log, "POST /v1/test", "Missing request entry")
	assert.Contains(t, log, "agent=curl", "Missing user agent")
	assert.Contains(t, log, "ip=210.172.60.240", "Missing user IP")
}
