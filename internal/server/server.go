package server

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/vliubezny/gstore/internal/service"
)

type server struct {
	s service.Service
}

// SetupRouter setups routes and handlers.
func SetupRouter(s service.Service, r chi.Router, username, password string) {
	srv := &server{
		s: s,
	}

	r.Use(
		loggerMiddleware,
		setContentTypeMiddleware(contentTypeJSON),
		recoveryMiddleware,
	)

	r.Get("/v1/categories", srv.getCategoriesHandler)
	r.Get("/v1/categories/{id}", srv.getCategoryHandler)
	r.Get("/v1/categories/{id}/products", srv.getCategoryProductsHandler)

	r.Get("/v1/stores", srv.getStoresHandler)
	r.Get("/v1/stores/{id}", srv.getStoreHandler)
	r.Get("/v1/stores/{id}/positions", srv.getStorePositionsHandler)

	r.Get("/v1/products/{id}/offers", srv.getProductOffersHandler)

	r.Group(func(r chi.Router) {
		r.Use(basicAuthMiddleware(username, password))

		r.Post("/v1/categories", srv.createCategoryHandler)
		r.Put("/v1/categories/{id}", srv.updateCategoryHandler)
		r.Delete("/v1/categories/{id}", srv.deleteCategoryHandler)

		r.Post("/v1/stores", srv.createStoreHandler)
		r.Put("/v1/stores/{id}", srv.updateStoreHandler)
		r.Delete("/v1/stores/{id}", srv.deleteStoreHandler)
	})

	decimal.MarshalJSONWithoutQuotes = true
}

func getLogger(r *http.Request) logrus.FieldLogger {
	return r.Context().Value(loggerKey{}).(logrus.FieldLogger)
}

func writeError(l logrus.FieldLogger, w http.ResponseWriter, code int, message string) {
	l.Error(message)

	body, _ := json.Marshal(errorResponse{
		Error: message,
	})

	w.WriteHeader(code)
	w.Write(body)
}

func writeInternalError(l logrus.FieldLogger, w http.ResponseWriter, message string) {
	l.Errorf("%s\n%s", message, string(debug.Stack()))

	body, _ := json.Marshal(errorResponse{
		Error: "internal error",
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(body)
}

func writeOK(l logrus.FieldLogger, w http.ResponseWriter, payload interface{}) {
	body, err := json.Marshal(payload)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to serialize payload")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
