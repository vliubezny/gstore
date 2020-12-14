package server

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func (s *server) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := s.s.GetCategories(r.Context())
	if err != nil {
		writeInternalError(w, err)
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

	body, err := json.Marshal(resp)
	if err != nil {
		writeInternalError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func writeInternalError(w http.ResponseWriter, err error) {
	logrus.WithError(err).Error("internal error")
	body, _ := json.Marshal(Error{
		Error: "internal error",
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(body)
}
