package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/vliubezny/gstore/internal/storage"
)

type extContext interface {
	sqlx.ExtContext
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type pg struct {
	dbx *sqlx.DB
	ext extContext
}

// New creates postgres storage.
func New(db *sql.DB) storage.Storage {
	dbx := sqlx.NewDb(db, "postgres")
	return pg{
		dbx: dbx,
		ext: dbx,
	}
}
