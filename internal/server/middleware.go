package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
	"github.com/vliubezny/gstore/internal/auth"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

type loggerKey struct{}

type claimsKey struct{}

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

// jwtAuthMiddleware authenticates user with JWT.
func jwtAuthMiddleware(accessTokenValidator auth.AccessTokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l := getLogger(r)

			token := extractBearer(r)
			if token == "" {
				writeError(l, w, http.StatusUnauthorized, "missing token")
				return
			}

			claims, err := accessTokenValidator(token)
			if err != nil {
				if errors.Is(err, auth.ErrInvalidToken) {
					writeError(l.WithError(err), w, http.StatusUnauthorized, "invalid access token")
					return
				}

				writeInternalError(l.WithError(err), w, "failed to validate access token")
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey{}, claims)
			ctx = context.WithValue(ctx, loggerKey{}, l.WithField("userID", claims.UserID))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// allowAdminMiddleware authorizes admin to access resource.
func allowAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := getLogger(r)

		claims, ok := r.Context().Value(claimsKey{}).(auth.AccessTokenClaims)
		if !ok {
			writeError(l, w, http.StatusUnauthorized, "authentication required")
			return
		}

		if !claims.IsAdmin {
			writeError(l, w, http.StatusForbidden, "access not allowed")
			return
		}

		next.ServeHTTP(w, r)
	})
}
