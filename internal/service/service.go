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
)

// Service provides business logic methods.
type Service interface {
	// GetCategories returns slice of product categories.
	GetCategories(ctx context.Context) ([]model.Category, error)

	// GetCategory returns a product category by ID.
	GetCategory(ctx context.Context, categoryID int64) (model.Category, error)

	// CreateCategory creates new category.
	CreateCategory(ctx context.Context, category model.Category) (model.Category, error)

	// UpdateCategory updates new category.
	UpdateCategory(ctx context.Context, category model.Category) error

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

func (s *service) GetStoreItems(ctx context.Context, storeID int64) ([]*model.Item, error) {
	stores, err := s.s.GetStoreItems(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store items: %w", err)
	}
	return stores, nil
}

func (s *service) GetStoreItem(ctx context.Context, itemID int64) (*model.Item, error) {
	item, err := s.s.GetStoreItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get store item: %w", err)
	}
	return item, nil
}

func (s *service) CreateStoreItem(ctx context.Context, item *model.Item) error {
	if err := s.s.CreateStoreItem(ctx, item); err != nil {
		return fmt.Errorf("failed to create store item: %w", err)
	}
	return nil
}

func (s *service) UpdateStoreItem(ctx context.Context, item *model.Item) error {
	if err := s.s.UpdateStoreItem(ctx, item); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update store item: %w", err)
	}
	return nil
}

func (s *service) DeleteStoreItem(ctx context.Context, itemID int64) error {
	if err := s.s.DeleteStoreItem(ctx, itemID); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete store item: %w", err)
	}
	return nil
}
