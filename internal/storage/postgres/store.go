package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (p pg) GetStores(ctx context.Context) ([]model.Store, error) {
	var stores []store
	if err := p.ext.SelectContext(ctx, &stores, "SELECT id, name FROM store"); err != nil {
		return nil, err
	}

	data := make([]model.Store, len(stores))
	for i, c := range stores {
		data[i] = c.toModel()
	}

	return data, nil
}

func (p pg) GetStore(ctx context.Context, storeID int64) (model.Store, error) {
	var s store
	err := p.ext.GetContext(ctx, &s, "SELECT id, name FROM store WHERE id = $1", storeID)

	if err == sql.ErrNoRows {
		return model.Store{}, storage.ErrNotFound
	}

	if err != nil {
		return model.Store{}, fmt.Errorf("failed to get store: %w", err)
	}

	return s.toModel(), nil
}

func (p pg) CreateStore(ctx context.Context, store model.Store) (model.Store, error) {
	if err := p.ext.GetContext(ctx, &store.ID, "INSERT INTO store (name) VALUES ($1) RETURNING id", store.Name); err != nil {
		return model.Store{}, fmt.Errorf("failed to create store: %w", err)
	}
	return store, nil
}

func (p pg) UpdateStore(ctx context.Context, store model.Store) error {
	res, err := p.ext.ExecContext(ctx, "UPDATE store SET name = $1 WHERE id = $2", store.Name, store.ID)

	if err != nil {
		return fmt.Errorf("failed to update store: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) DeleteStore(ctx context.Context, storeID int64) error {
	res, err := p.ext.ExecContext(ctx, "DELETE FROM store WHERE id = $1", storeID)

	if err != nil {
		return fmt.Errorf("failed to delete store: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}
