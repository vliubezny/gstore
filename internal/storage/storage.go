package storage

import (
	"context"
	"errors"
	"time"

	"github.com/vliubezny/gstore/internal/model"
)

//go:generate mockgen -destination=./storage_mock.go -package=storage -source=storage.go

var (
	// ErrNotFound states that record was not found.
	ErrNotFound = errors.New("not found")

	// ErrUnknownCategory states that category is unknown.
	ErrUnknownCategory = errors.New("category is unknown")

	// ErrUnknownStore states that store is unknown.
	ErrUnknownStore = errors.New("store is unknown")

	// ErrUnknownProduct states that product is unknown.
	ErrUnknownProduct = errors.New("product is unknown")

	// ErrEmailIsTaken states that email address is taken.
	ErrEmailIsTaken = errors.New("email is taken")
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
	GetStores(ctx context.Context) ([]model.Store, error)

	// GetStore returns a product store by ID.
	GetStore(ctx context.Context, storeID int64) (model.Store, error)

	// CreateStore creates new store.
	CreateStore(ctx context.Context, store model.Store) (model.Store, error)

	// UpdateStore updates store.
	UpdateStore(ctx context.Context, store model.Store) error

	// DeleteStore deletes store from storage.
	DeleteStore(ctx context.Context, storeID int64) error

	// GetProducts returns slice of products in category.
	GetProducts(ctx context.Context, categoryID int64) ([]model.Product, error)

	// GetProduct returns a product by ID.
	GetProduct(ctx context.Context, productID int64) (model.Product, error)

	// CreateProduct creates new product.
	CreateProduct(ctx context.Context, product model.Product) (model.Product, error)

	// UpdateProduct updates product.
	UpdateProduct(ctx context.Context, product model.Product) error

	// DeleteProduct deletes product.
	DeleteProduct(ctx context.Context, productID int64) error

	// GetStorePositions returns slice of store positions.
	GetStorePositions(ctx context.Context, storeID int64) ([]model.Position, error)

	// GetProductPositions returns slice of product positions.
	GetProductPositions(ctx context.Context, productID int64) ([]model.Position, error)

	// UpsertPosition updates position or creates new one if it doesn't exist.
	UpsertPosition(ctx context.Context, position model.Position) error

	// DeletePosition deletes position.
	DeletePosition(ctx context.Context, productID, storeID int64) error
}

// UserStorage provides methods to interact with user storage.
type UserStorage interface {
	// CreateUser creates new user.
	CreateUser(ctx context.Context, user model.User) (model.User, error)

	// GetUserByEmail returns user from storage by email.
	GetUserByEmail(ctx context.Context, email string) (model.User, error)

	// SaveToken saves token reference.
	SaveToken(ctx context.Context, tokenID string, userID int64, expiresAt time.Time) error

	// DeleteToken deletes token.
	DeleteToken(ctx context.Context, tokenID string) error
}
