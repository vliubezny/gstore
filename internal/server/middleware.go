package server

import "net/http"

// setContentTypeMiddleware sets default content type.
func setContentTypeMiddleware(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(headerContentType, contentTypeJSON)
			next.ServeHTTP(w, r)
		})
	}
}
