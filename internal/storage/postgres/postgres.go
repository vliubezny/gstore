package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

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
