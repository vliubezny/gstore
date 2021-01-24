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
	GetCategories(ctx context.Context) ([]model.Category, error)

	// GetCategory returns a product category by ID.
	GetCategory(ctx context.Context, categoryID int64) (model.Category, error)

	// CreateCategory creates new category.
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)

	// UpdateCategory updates category.
	UpdateCategory(ctx context.Context, category model.Category) error

	// DeleteCategory deletes category from storage.
	DeleteCategory(ctx context.Context, categoryID int64) error

	// GetStores returns slice of stores.
	GetStores(ctx context.Context) ([]*model.Store, error)

	// GetStore returns a product store by ID.
	GetStore(ctx context.Context, storeID int64) (*model.Store, error)

	// CreateStore creates new store.
	CreateStore(ctx context.Context, store *model.Store) error

	// UpdateStore updates store.
	UpdateStore(ctx context.Context, store *model.Store) error

	// DeleteStore deletes store from storage.
	DeleteStore(ctx context.Context, storeID int64) error

	// GetProducts returns slice of products in category.
	GetProducts(ctx context.Context, categoryID int64) ([]*model.Product, error)

	// GetProduct returns a product by ID.
	GetProduct(ctx context.Context, productID int64) (*model.Product, error)

	// CreateProduct creates new product.
	CreateProduct(ctx context.Context, product *model.Product) error

	// UpdateProduct updates product.
	UpdateProduct(ctx context.Context, product *model.Product) error

	// DeleteProduct deletes product.
	DeleteProduct(ctx context.Context, productID int64) error
}
