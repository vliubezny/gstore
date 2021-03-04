package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

const (
	productIDFKConstraint = "position_product_id_fkey"
	storeIDFKConstraint   = "position_store_id_fkey"
)

func (p pg) GetStorePositions(ctx context.Context, storeID int64) ([]model.Position, error) {
	var positions []position

	if err := p.db.SelectContext(ctx, &positions, "SELECT product_id, store_id, price FROM position WHERE store_id=$1", storeID); err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	data := make([]model.Position, len(positions))
	for i, d := range positions {
		data[i] = d.toModel()
	}

	return data, nil
}

func (p pg) GetProductPositions(ctx context.Context, productID int64) ([]model.Position, error) {
	var positions []position

	if err := p.db.SelectContext(ctx, &positions, "SELECT product_id, store_id, price FROM position WHERE product_id=$1", productID); err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	data := make([]model.Position, len(positions))
	for i, d := range positions {
		data[i] = d.toModel()
	}

	return data, nil
}

func (p pg) UpsertPosition(ctx context.Context, position model.Position) error {
	if _, err := p.db.ExecContext(ctx, `
		INSERT INTO position (product_id, store_id, price) VALUES($1, $2, $3)
			ON CONFLICT(product_id, store_id) DO UPDATE SET price = EXCLUDED.price;
	`, position.ProductID, position.StoreID, position.Price); err != nil {
		if err, ok := err.(*pq.Error); ok {
			switch err.Constraint {
			case productIDFKConstraint:
				return storage.ErrUnknownProduct
			case storeIDFKConstraint:
				return storage.ErrUnknownStore
			}
		}

		return fmt.Errorf("failed to upsert position: %w", err)
	}

	return nil
}

func (p pg) DeletePosition(ctx context.Context, productID, storeID int64) error {
	res, err := p.db.ExecContext(ctx, `
		DELETE FROM position WHERE product_id = $1 AND store_id = $2
	`, productID, storeID)

	if err != nil {
		return fmt.Errorf("failed to delete position: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}
