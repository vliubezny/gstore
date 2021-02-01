package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

type loggerKey struct{}

// setContentTypeMiddleware sets default content type.
func setContentTypeMiddleware(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(headerContentType, contentTypeJSON)
			next.ServeHTTP(w, r)
		})
	}
}

// loggerMiddleware populates request context with logger and logs request entry.
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logrus.WithFields(logrus.Fields{
			"ip":    realip.FromRequest(r),
			"agent": r.UserAgent(),
		})
		ctx := context.WithValue(r.Context(), loggerKey{}, logger)
		logger.Debugf("%s %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// recoveryMiddleware recovers after panic.
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				writeInternalError(getLogger(r), w, fmt.Sprintf("recover from panic: %s\n", spew.Sdump(e)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// basicAuthMiddleware handles basic authentication.
func basicAuthMiddleware(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := getLogger(r)

			u, p, ok := r.BasicAuth()
			if !ok {
				writeError(l, w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			if u != username || p != password {
				writeError(l.WithField("username", u), w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
