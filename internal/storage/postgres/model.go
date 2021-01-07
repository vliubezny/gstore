package postgres

import "github.com/vliubezny/gstore/internal/model"

type category struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (c category) toModel() *model.Category {
	return &model.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

type store struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func (s store) toModel() *model.Store {
	return &model.Store{
		ID:   s.ID,
		Name: s.Name,
	}
}

type item struct {
	ID          int64  `db:"id"`
	StoreID     int64  `db:"store_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Price       int64  `db:"price"`
}

func (i item) toModel() *model.Item {
	return &model.Item{
		ID:          i.ID,
		StoreID:     i.StoreID,
		Name:        i.Name,
		Description: i.Description,
		Price:       i.Price,
	}
}
