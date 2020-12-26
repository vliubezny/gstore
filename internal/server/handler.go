package server

import (
	"net/http"
	"strconv"
)

func (s *server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categories, err := s.s.GetCategories(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get categories")
		return
	}

	resp := getCategoriesResponse{
		Categories: make([]*category, len(categories)),
	}

	for i, c := range categories {
		resp.Categories[i] = &category{
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

	resp := getStoresResponse{
		Stores: make([]*store, len(stores)),
	}

	for i, c := range stores {
		resp.Stores[i] = &store{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	writeOK(l, w, resp)
}

func (s *server) getStoreItemsHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	sid := r.URL.Query().Get("storeId")
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

	resp := getStoreItemsResponse{
		Items: make([]*item, len(items)),
	}

	for i, c := range items {
		resp.Items[i] = &item{
			ID:          c.ID,
			StoreID:     c.StoreID,
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
		}
	}

	writeOK(l, w, resp)
}
