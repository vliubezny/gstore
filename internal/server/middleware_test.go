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

func Test_basicAuthMiddleware(t *testing.T) {
	testCases := []struct {
		desc     string
		username string
		password string
		rcode    int
		rdata    string
	}{
		{
			desc:     "allow valid credentials",
			username: testUsername,
			password: testPassword,
			rcode:    http.StatusOK,
			rdata:    `{"result":"OK"}`,
		},
		{
			desc:     "block anonymous",
			username: "",
			password: "",
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
		{
			desc:     "block invalid username",
			username: "invalid",
			password: testPassword,
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
		{
			desc:     "block invalid password",
			username: testUsername,
			password: "invalid",
			rcode:    http.StatusUnauthorized,
			rdata:    `{"error":"Unauthorized"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			logger, _ := test.NewNullLogger()
			ctx := context.WithValue(context.Background(), loggerKey{}, logger)
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", nil).WithContext(ctx)

			if tC.username != "" && tC.password != "" {
				req.SetBasicAuth(tC.username, tC.password)
			}

			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"result":"OK"}`))
			})

			basicAuthMiddleware(testUsername, testPassword)(h).ServeHTTP(rec, req)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}
