package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (p pg) GetStores(ctx context.Context) ([]*model.Store, error) {
	var stores []store
	if err := p.db.SelectContext(ctx, &stores, "SELECT * FROM store"); err != nil {
		return nil, err
	}

	data := make([]*model.Store, len(stores))
	for i, c := range stores {
		data[i] = c.toModel()
	}

	return data, nil
}

func (p pg) GetStore(ctx context.Context, storeID int64) (*model.Store, error) {
	var s store
	err := p.db.GetContext(ctx, &s, "SELECT * FROM store WHERE id = $1", storeID)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return s.toModel(), nil
}

func (p pg) CreateStore(ctx context.Context, store *model.Store) error {
	var id int64
	if err := p.db.GetContext(ctx, &id, "INSERT INTO store (name) VALUES ($1) RETURNING id", store.Name); err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	store.ID = id
	return nil
}

func (p pg) UpdateStore(ctx context.Context, store *model.Store) error {
	res, err := p.db.ExecContext(ctx, "UPDATE store SET name = $1 WHERE id = $2", store.Name, store.ID)

	if err != nil {
		return fmt.Errorf("failed to update store: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) DeleteStore(ctx context.Context, storeID int64) error {
	res, err := p.db.ExecContext(ctx, "DELETE FROM store WHERE id = $1", storeID)

	if err != nil {
		return fmt.Errorf("failed to delete store: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}