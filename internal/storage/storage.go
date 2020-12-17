package storage

import (
	"context"

	"github.com/vliubezny/gstore/internal/model"
)

//go:generate mockgen -destination=./storage_mock.go -package=storage -source=storage.go

// Storage provides methods to interact with data storage.
type Storage interface {
	// GetCategories returns slice of product categories.
	GetCategories(ctx context.Context) ([]*model.Category, error)
}
