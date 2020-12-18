package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"github.com/vliubezny/gstore/internal/service"
)

type server struct {
	s service.Service
}

// SetupRouter setups routes and handlers.
func SetupRouter(s service.Service, r chi.Router) {
	srv := &server{
		s: s,
	}

	r.Use(
		loggerMiddleware,
		setContentTypeMiddleware(contentTypeJSON),
	)

	r.Get("/v1/categories", srv.getCategoriesHandler)
}

func getLogger(r *http.Request) logrus.FieldLogger {
	return r.Context().Value(loggerKey{}).(logrus.FieldLogger)
}
