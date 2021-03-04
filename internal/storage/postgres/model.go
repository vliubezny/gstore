package postgres

import (
	"github.com/shopspring/decimal"
	"github.com/vliubezny/gstore/internal/model"
)

type category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (c category) toModel() model.Category {
	return model.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

type store struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s store) toModel() model.Store {
	return model.Store{
		ID:   s.ID,
		Name: s.Name,
	}
}

type product struct {
	ID          int64  `db:"id"`
	CategoryID  int64  `db:"category_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func (i product) toModel() model.Product {
	return model.Product{
		ID:          i.ID,
		CategoryID:  i.CategoryID,
		Name:        i.Name,
		Description: i.Description,
	}
}

type position struct {
	ProductID int64           `db:"product_id"`
	StoreID   int64           `db:"store_id"`
	Price     decimal.Decimal `db:"price"`
}

func (p position) toModel() model.Position {
	return model.Position{
		ProductID: p.ProductID,
		StoreID:   p.StoreID,
		Price:     p.Price,
	}
}

type user struct {
	ID           int64  `db:"id"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
	IsAdmin      bool   `db:"is_admin"`
}

func (u user) toModel() model.User {
	return model.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		IsAdmin:      u.IsAdmin,
	}
}
