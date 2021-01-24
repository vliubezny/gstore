package server

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/vliubezny/gstore/internal/model"
)

// errorResponse represents error response
type errorResponse struct {
	Error string `json:"error"`
}

// category represents category object.
type category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func fromCategoryModel(c model.Category) category {
	return category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c category) toModel() model.Category {
	return model.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c category) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(2, 80)),
	)
}

type store struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func newStore(s *model.Store) store {
	return store{
		ID:   s.ID,
		Name: s.Name,
	}
}

func (s store) toModel() *model.Store {
	return &model.Store{
		ID:   s.ID,
		Name: s.Name,
	}
}

func (s store) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(2, 80)),
	)
}

type item struct {
	ID          int64  `json:"id"`
	StoreID     int64  `json:"storeId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}
