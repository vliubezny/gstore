package server

import (
	"net/http"
)

func (s *server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	l := getLogger(r)

	categories, err := s.s.GetCategories(r.Context())
	if err != nil {
		writeInternalError(l.WithError(err), w, "fail to get categories")
		return
	}

	resp := GetCategoriesResponse{
		Categories: make([]*Category, len(categories)),
	}

	for i, c := range categories {
		resp.Categories[i] = &Category{
			ID:   c.ID,
			Name: c.Name,
		}
	}

	writeOK(l, w, resp)
}
