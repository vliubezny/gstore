package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (p pg) GetStoreItems(ctx context.Context, storeID int64) ([]*model.Item, error) {
	var items []item

	if err := p.db.SelectContext(ctx, &items, "SELECT * FROM item WHERE store_id=$1", storeID); err != nil {
		return nil, err
	}

	data := make([]*model.Item, len(items))
	for i, d := range items {
		data[i] = d.toModel()
	}

	return data, nil
}

func (p pg) GetStoreItem(ctx context.Context, itemID int64) (*model.Item, error) {
	var itm item
	err := p.db.GetContext(ctx, &itm, "SELECT * FROM item WHERE id = $1", itemID)

	if err == sql.ErrNoRows {
		return nil, storage.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get store item: %w", err)
	}

	return itm.toModel(), nil
}

func (p pg) CreateStoreItem(ctx context.Context, item *model.Item) error {
	var id int64
	if err := p.db.GetContext(ctx, &id, "INSERT INTO item (store_id, name, description, price) VALUES ($1, $2, $3, $4) RETURNING id",
		item.StoreID, item.Name, item.Description, item.Price); err != nil {
		return fmt.Errorf("failed to create store item: %w", err)
	}

	item.ID = id
	return nil
}

func (p pg) UpdateStoreItem(ctx context.Context, item *model.Item) error {
	res, err := p.db.ExecContext(ctx, `
		UPDATE item SET
		store_id =$2,
		name = $3,
		description = $4,
		price = $5
		WHERE id = $1
	`, item.ID, item.StoreID, item.Name, item.Description, item.Price)

	if err != nil {
		return fmt.Errorf("failed to update store item: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) DeleteStoreItem(ctx context.Context, itemID int64) error {
	res, err := p.db.ExecContext(ctx, "DELETE FROM item WHERE id = $1", itemID)

	if err != nil {
		return fmt.Errorf("failed to delete store item: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}
