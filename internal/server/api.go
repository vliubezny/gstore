package server

import (
	"github.com/shopspring/decimal"
	"github.com/vliubezny/gstore/internal/model"
)

// errorResponse represents error response
type errorResponse struct {
	Error string `json:"error"`
}

// category represents category object.
type category struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required,gte=2,lte=80"`
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

type store struct {
	ID   int64  `json:"id"`
	Name string `json:"name" validate:"required,gte=2,lte=80"`
}

func fromStoreModel(s model.Store) store {
	return store{
		ID:   s.ID,
		Name: s.Name,
	}
}

func (s store) toModel() model.Store {
	return model.Store{
		ID:   s.ID,
		Name: s.Name,
	}
}

type product struct {
	ID          int64  `json:"id"`
	CategoryID  int64  `json:"categoryId" validate:"required"`
	Name        string `json:"name" validate:"required,gte=3,lte=160"`
	Description string `json:"description" validate:"required"`
}

func fromProductModel(p model.Product) product {
	return product{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
	}
}

func (p product) toModel() model.Product {
	return model.Product{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
	}
}

type position struct {
	ProductID int64           `json:"productId"`
	StoreID   int64           `json:"storeId"`
	Price     decimal.Decimal `json:"price" validate:"gt=0"`
}

func fromPositionModel(p model.Position) position {
	return position{
		ProductID: p.ProductID,
		StoreID:   p.StoreID,
		Price:     p.Price,
	}
}

func (p position) toModel() model.Position {
	return model.Position{
		ProductID: p.ProductID,
		StoreID:   p.StoreID,
		Price:     p.Price,
	}
}
