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
