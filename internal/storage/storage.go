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

	// GetCategory returns a product category by ID.
	GetCategory(ctx context.Context, categoryID int64) (*model.Category, error)

	// CreateCategory creates new category.
	CreateCategory(ctx context.Context, category *model.Category) error

	// UpdateCategory updates new category.
	UpdateCategory(ctx context.Context, category *model.Category) error

	// DeleteCategory deletes category from storage.
	DeleteCategory(ctx context.Context, categoryID int64) error

	// GetStores returns slice of stores.
	GetStores(ctx context.Context) ([]*model.Store, error)

	// GetStore returns a product store by ID.
	GetStore(ctx context.Context, storeID int64) (*model.Store, error)

	// CreateStore creates new store.
	CreateStore(ctx context.Context, store *model.Store) error

	// UpdateStore updates new store.
	UpdateStore(ctx context.Context, store *model.Store) error

	// DeleteStore deletes store from storage.
	DeleteStore(ctx context.Context, storeID int64) error

	// GetStoreItems returns slice of store items.
	GetStoreItems(ctx context.Context, storeID int64) ([]*model.Item, error)

	// GetStoreItem returns a store item by ID.
	GetStoreItem(ctx context.Context, itemID int64) (*model.Item, error)

	// CreateStoreItem creates new store item.
	CreateStoreItem(ctx context.Context, item *model.Item) error

	// UpdateStoreItem updates new store item.
	UpdateStoreItem(ctx context.Context, item *model.Item) error

	// DeleteStoreItem deletes store item from storage.
	DeleteStoreItem(ctx context.Context, itemID int64) error
}
