package mem

import (
	"context"

	"github.com/vliubezny/gstore/internal/storage"
)

type memStorage struct {
	s []*storage.Category
}

// New creates prepopulated in-memory storage.
func New() storage.Storage {
	s := &memStorage{}
	s.s = []*storage.Category{
		{ID: 1, Name: "Electronics"},
		{ID: 2, Name: "Computers"},
		{ID: 3, Name: "Smart Home"},
		{ID: 4, Name: "Arts & Crafts"},
		{ID: 5, Name: "Health & Household"},
		{ID: 6, Name: "Automotive"},
		{ID: 7, Name: "Pet supplies"},
		{ID: 8, Name: "Software"},
		{ID: 9, Name: "Sports & Outdoors"},
		{ID: 10, Name: "Toys and Games"},
	}
	return s
}

func (s *memStorage) GetCategories(ctx context.Context) ([]*storage.Category, error) {
	return s.s, nil
}
