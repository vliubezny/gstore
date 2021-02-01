package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

//go:generate mockgen -destination=./service_mock.go -package=service -source=service.go

var (
	// ErrNotFound states that object(s) was not found.
	ErrNotFound = errors.New("not found")

	// ErrUnknownStore states that store is unknown.
	ErrUnknownStore = errors.New("store is unknown")

	// ErrUnknownProduct states that product is unknown.
	ErrUnknownProduct = errors.New("product is unknown")
)

// Service provides business logic methods.
type Service interface {
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

	// GetStorePositions returns slice of store positions.
	GetStorePositions(ctx context.Context, storeID int64) ([]model.Position, error)

	// GetProductPositions returns slice of product positions.
	GetProductPositions(ctx context.Context, productID int64) ([]model.Position, error)

	// SetPosition updates position or creates new one if it doesn't exist.
	SetPosition(ctx context.Context, position model.Position) error

	// DeletePosition deletes position.
	DeletePosition(ctx context.Context, productID, storeID int64) error
}

type service struct {
	s storage.Storage
}

// New creates service instance.
func New(s storage.Storage) Service {
	return &service{
		s: s,
	}
}

func (s *service) GetCategories(ctx context.Context) ([]model.Category, error) {
	categories, err := s.s.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (s *service) GetCategory(ctx context.Context, categoryID int64) (model.Category, error) {
	category, err := s.s.GetCategory(ctx, categoryID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return model.Category{}, ErrNotFound
		}
		return model.Category{}, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}

func (s *service) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	c, err := s.s.CreateCategory(ctx, category)
	if err != nil {
		return model.Category{}, fmt.Errorf("failed to create category: %w", err)
	}
	return c, nil
}

func (s *service) UpdateCategory(ctx context.Context, category model.Category) error {
	if err := s.s.UpdateCategory(ctx, category); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

func (s *service) DeleteCategory(ctx context.Context, categoryID int64) error {
	if err := s.s.DeleteCategory(ctx, categoryID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (s *service) GetStores(ctx context.Context) ([]*model.Store, error) {
	stores, err := s.s.GetStores(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stores: %w", err)
	}
	return stores, nil
}

func (s *service) GetStore(ctx context.Context, storeID int64) (*model.Store, error) {
	store, err := s.s.GetStore(ctx, storeID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get store: %w", err)
	}
	return store, nil
}

func (s *service) CreateStore(ctx context.Context, store *model.Store) error {
	if err := s.s.CreateStore(ctx, store); err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	return nil
}

func (s *service) UpdateStore(ctx context.Context, store *model.Store) error {
	if err := s.s.UpdateStore(ctx, store); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update store: %w", err)
	}
	return nil
}

func (s *service) DeleteStore(ctx context.Context, storeID int64) error {
	if err := s.s.DeleteStore(ctx, storeID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete store: %w", err)
	}
	return nil
}

func (s *service) GetProducts(ctx context.Context, categoryID int64) ([]*model.Product, error) {
	products, err := s.s.GetProducts(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	return products, nil
}

func (s *service) GetProduct(ctx context.Context, productID int64) (*model.Product, error) {
	product, err := s.s.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *service) CreateProduct(ctx context.Context, product *model.Product) error {
	if err := s.s.CreateProduct(ctx, product); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (s *service) UpdateProduct(ctx context.Context, product *model.Product) error {
	if err := s.s.UpdateProduct(ctx, product); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

func (s *service) DeleteProduct(ctx context.Context, productID int64) error {
	if err := s.s.DeleteProduct(ctx, productID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

func (s *service) GetStorePositions(ctx context.Context, storeID int64) ([]model.Position, error) {
	positions, err := s.s.GetStorePositions(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store positions: %w", err)
	}
	return positions, nil
}

func (s *service) GetProductPositions(ctx context.Context, productID int64) ([]model.Position, error) {
	positions, err := s.s.GetProductPositions(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product positions: %w", err)
	}
	return positions, nil
}

func (s *service) SetPosition(ctx context.Context, position model.Position) error {
	if err := s.s.UpsertPosition(ctx, position); err != nil {
		switch {
		case errors.Is(err, storage.ErrUnknownProduct):
			return ErrUnknownProduct
		case errors.Is(err, storage.ErrUnknownStore):
			return ErrUnknownStore
		}
		return fmt.Errorf("failed to set position: %w", err)
	}

	return nil
}

func (s *service) DeletePosition(ctx context.Context, productID, storeID int64) error {
	if err := s.s.DeletePosition(ctx, productID, storeID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete position: %w", err)
	}

	return nil
}
