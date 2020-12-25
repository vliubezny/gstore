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
	GetCategories(ctx context.Context) ([]*model.Category, error)

	// GetStoreItems returns slice of store items.
	GetStoreItems(ctx context.Context, storeID int64) ([]*model.Item, error)

	// GetStores returns slice of stores.
	GetStores(ctx context.Context) ([]*model.Store, error)
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

func (s *service) GetCategories(ctx context.Context) ([]*model.Category, error) {
	return s.s.GetCategories(ctx)
}

func (s *service) GetStoreItems(ctx context.Context, storeID int64) ([]*model.Item, error) {
	stores, err := s.s.GetStoreItems(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get store items: %w", err)
	}
	return stores, nil
}

func (s *service) GetStores(ctx context.Context) ([]*model.Store, error) {
	stores, err := s.s.GetStores(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stores: %w", err)
	}
	return stores, nil
}
