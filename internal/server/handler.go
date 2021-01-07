package server

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (s *server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categories, err := s.s.GetCategories(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get categories")
		return
	}

	resp := make([]*category, len(categories))

	for i, c := range categories {
		resp[i] = &category{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	writeOK(l, w, resp)
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
		writeError(l.WithError(err), w, http.StatusBadRequest, "invalid storeId")
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
