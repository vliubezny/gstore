package service

import (
	"context"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

//go:generate mockgen -destination=./service_mock.go -package=service -source=service.go

// Service provides business logic methods.
type Service interface {
	// GetCategories returns slice of product categories.
	GetCategories(ctx context.Context) ([]*model.Category, error)
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
