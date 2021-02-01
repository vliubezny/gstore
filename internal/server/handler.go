package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/vliubezny/gstore/internal/service"
)

func (s *server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categories, err := s.s.GetCategories(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get categories")
		return
	}

	resp := make([]category, len(categories))

	for i, c := range categories {
		resp[i] = fromCategoryModel(c)
	}

	writeOK(l, w, resp)
}

func (s *server) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	c, err := s.s.GetCategory(r.Context(), categoryID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "category not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to get category")
		return
	}

	writeOK(l, w, fromCategoryModel(c))
}

func (s *server) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	c, err := s.s.CreateCategory(r.Context(), req.toModel())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to create category")
		return
	}

	writeOK(l, w, fromCategoryModel(c))
}

func (s *server) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	var req category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	c := req.toModel()
	c.ID = categoryID

	if err := s.s.UpdateCategory(r.Context(), c); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "category not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to update category")
		return
	}

	writeOK(l, w, fromCategoryModel(c))
}

func (s *server) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	err = s.s.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "category not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to delete category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getStoresHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	stores, err := s.s.GetStores(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get stores")
		return
	}

	resp := make([]*store, len(stores))

	for i, c := range stores {
		resp[i] = &store{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	writeOK(l, w, resp)
}

func (s *server) getStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	storeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	str, err := s.s.GetStore(r.Context(), storeID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to get store")
		return
	}

	writeOK(l, w, fromStoreModel(str))
}

func (s *server) createStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req store
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	str := req.toModel()
	if err := s.s.CreateStore(r.Context(), str); err != nil {
		writeInternalError(l.WithError(err), w, "fail to create store")
		return
	}

	writeOK(l, w, fromStoreModel(str))
}

func (s *server) updateStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	storeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	var req store
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validate(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	str := req.toModel()
	str.ID = storeID

	if err := s.s.UpdateStore(r.Context(), str); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to update store")
		return
	}

	writeOK(l, w, fromStoreModel(str))
}

func (s *server) deleteStoreHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	storeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	err = s.s.DeleteStore(r.Context(), storeID)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(l.WithError(err), w, http.StatusNotFound, "store not found")
			return
		}

		writeInternalError(l.WithError(err), w, "fail to delete store")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *server) getStorePositionsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	storeID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	positions, err := s.s.GetStorePositions(r.Context(), storeID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get store positions")
		return
	}

	resp := make([]position, len(positions))

	for i, p := range positions {
		resp[i] = fromPositionModel(p)
	}

	writeOK(l, w, resp)
}

func (s *server) getCategoryProductsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	categoryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid category ID")
		return
	}

	products, err := s.s.GetProducts(r.Context(), categoryID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get products")
		return
	}

	resp := make([]*product, len(products))

	for i, c := range products {
		resp[i] = &product{
			ID:          c.ID,
			CategoryID:  c.CategoryID,
			Name:        c.Name,
			Description: c.Description,
		}
	}

	writeOK(l, w, resp)
}

func (s *server) getProductOffersHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	id := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid product ID")
		return
	}

	positions, err := s.s.GetProductPositions(r.Context(), productID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get product offers")
		return
	}

	resp := make([]position, len(positions))

	for i, p := range positions {
		resp[i] = fromPositionModel(p)
	}

	writeOK(l, w, resp)
}
