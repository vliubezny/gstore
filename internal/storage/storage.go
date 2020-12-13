package storage

import "context"

// Category represents product cetegory.
type Category struct {
	ID   int64
	Name string
}

// Storage provides methods to interact with data storage.
type Storage interface {
	// GetCategories returns slice of product categories.
	GetCategories(ctx context.Context) ([]*Category, error)
}
