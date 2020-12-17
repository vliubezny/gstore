package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_setContentTypeMiddleware(t *testing.T) {
	rec := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, "", nil)
	require.NoError(t, err)

	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"result\":42}"))
	})

	setContentTypeMiddleware(contentTypeJSON)(h).ServeHTTP(rec, r)
	assert.Equal(t, "application/json", rec.Result().Header.Get("Content-Type"))
}
