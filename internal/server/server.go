package server

import (
	"github.com/go-chi/chi"
	"github.com/vliubezny/gstore/internal/service"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
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
		setContentTypeMiddleware(contentTypeJSON),
	)

	r.Get("/v1/categories", srv.getCategoriesHandler)
}
