package server

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/vliubezny/gstore/internal/auth"
	"github.com/vliubezny/gstore/internal/service"
)

type server struct {
	s service.Service
	a auth.Service
}

// SetupRouter setups routes and handlers.
func SetupRouter(s service.Service, a auth.Service, r chi.Router, accessTokenValidator auth.AccessTokenValidator) {
	srv := &server{
		s: s,
		a: a,
	}

	r.Use(
		loggerMiddleware,
		setContentTypeMiddleware(contentTypeJSON),
		recoveryMiddleware,
	)

	r.Post("/v1/register", srv.registerHandler)
	r.Post("/v1/login", srv.loginHandler)
	r.Post("/v1/refresh", srv.refreshHandler)
	r.Post("/v1/revoke", srv.revokeHandler)

	r.Get("/v1/categories", srv.getCategoriesHandler)
	r.Get("/v1/categories/{id}", srv.getCategoryHandler)
	r.Get("/v1/categories/{id}/products", srv.getCategoryProductsHandler)

	r.Get("/v1/stores", srv.getStoresHandler)
	r.Get("/v1/stores/{id}", srv.getStoreHandler)
	r.Get("/v1/stores/{id}/positions", srv.getStorePositionsHandler)

	r.Get("/v1/products/{id}", srv.getProductHandler)
	r.Get("/v1/products/{id}/offers", srv.getProductOffersHandler)

	r.Group(func(r chi.Router) {
		r.Use(
			jwtAuthMiddleware(accessTokenValidator),
			allowAdminMiddleware,
		)

		r.Put("/v1/users/{id}/permissions", srv.updateUserPermissionsHandler)

		r.Post("/v1/categories", srv.createCategoryHandler)
		r.Put("/v1/categories/{id}", srv.updateCategoryHandler)
		r.Delete("/v1/categories/{id}", srv.deleteCategoryHandler)

		r.Post("/v1/products", srv.createProductHandler)
		r.Put("/v1/products/{id}", srv.updateProductHandler)
		r.Delete("/v1/products/{id}", srv.deleteProductHandler)

		r.Post("/v1/stores", srv.createStoreHandler)
		r.Put("/v1/stores/{id}", srv.updateStoreHandler)
		r.Delete("/v1/stores/{id}", srv.deleteStoreHandler)

		r.Put("/v1/stores/{id}/positions/{productId}", srv.setPositionHandler)
		r.Delete("/v1/stores/{id}/positions/{productId}", srv.deletePositionHandler)
	})

	decimal.MarshalJSONWithoutQuotes = true
}

func getLogger(r *http.Request) logrus.FieldLogger {
	return r.Context().Value(loggerKey{}).(logrus.FieldLogger)
}

func extractBearer(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && strings.ToUpper(auth[0:7]) == "BEARER " {
		return auth[7:]
	}
	return ""
}

func getIDFromURL(r *http.Request, key string) (int64, error) {
	id := chi.URLParam(r, key)
	return strconv.ParseInt(id, 10, 64)
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
