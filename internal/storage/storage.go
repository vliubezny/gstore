package storage

import (
	"context"
	"errors"

	"github.com/vliubezny/gstore/internal/model"
)

//go:generate mockgen -destination=./storage_mock.go -package=storage -source=storage.go

var (
	// ErrNotFound states that record was not found.
	ErrNotFound = errors.New("not found")
)

// Storage provides methods to interact with data storage.
type Storage interface {
	// GetCategories returns slice of product categories.
	GetCategories(ctx context.Context) ([]*model.Category, error)

	// GetStoreItems returns slice of store items.
	GetStoreItems(ctx context.Context, storeID int64) ([]*model.Item, error)

	// GetStores returns slice of stores.
	GetStores(ctx context.Context) ([]*model.Store, error)
}
