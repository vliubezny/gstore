package server

import (
	"github.com/shopspring/decimal"
	"github.com/vliubezny/gstore/internal/auth"
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

type credentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8,lte=160"`
}

type user struct {
	ID      int64  `json:"id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

func fromUserModel(u model.User) user {
	return user{
		ID:      u.ID,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}

type tokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func fromTokenPairModel(tp auth.TokenPair) tokenPair {
	return tokenPair{
		AccessToken:  tp.AccessToken,
		RefreshToken: tp.RefreshToken,
	}
}

type userPermissions struct {
	IsAdmin bool `json:"isAdmin" validate:"required"`
}

func (p userPermissions) toModel() model.User {
	return model.User{
		IsAdmin: p.IsAdmin,
	}
}
