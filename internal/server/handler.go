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
		resp[i] = newCategory(c)
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

	writeOK(l, w, newCategory(c))
}

func (s *server) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	var req category
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, err.Error())
		return
	}

	c := req.toModel()
	if err := s.s.CreateCategory(r.Context(), c); err != nil {
		writeInternalError(l.WithError(err), w, "fail to create category")
		return
	}

	writeOK(l, w, newCategory(c))
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

	if err := req.Validate(); err != nil {
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

	writeOK(l, w, newCategory(c))
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

func (s *server) getStoreItemsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	sid := chi.URLParam(r, "id")
	storeID, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid store ID")
		return
	}

	items, err := s.s.GetStoreItems(r.Context(), storeID)
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get items")
		return
	}

	resp := make([]*item, len(items))

	for i, c := range items {
		resp[i] = &item{
			ID:          c.ID,
			StoreID:     c.StoreID,
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
		}
	}

	writeOK(l, w, resp)
}
