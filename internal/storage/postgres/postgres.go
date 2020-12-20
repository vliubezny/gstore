package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

type pg struct {
	db *sqlx.DB
}

// New creates postgres storage.
func New(db *sql.DB) storage.Storage {
	return pg{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (p pg) GetCategories(ctx context.Context) ([]*model.Category, error) {
	var categories []category
	if err := p.db.SelectContext(ctx, &categories, "SELECT * FROM category"); err != nil {
		return nil, err
	}

	data := make([]*model.Category, len(categories))
	for i, c := range categories {
		data[i] = c.toModel()
	}

	return data, nil
}
