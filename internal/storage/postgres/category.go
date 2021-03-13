package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (p pg) GetCategories(ctx context.Context) ([]model.Category, error) {
	var categories []category
	if err := p.ext.SelectContext(ctx, &categories, "SELECT id, name FROM category"); err != nil {
		return nil, err
	}

	data := make([]model.Category, len(categories))
	for i, c := range categories {
		data[i] = c.toModel()
	}

	return data, nil
}

func (p pg) GetCategory(ctx context.Context, categoryID int64) (model.Category, error) {
	var c category
	err := p.ext.GetContext(ctx, &c, "SELECT id, name FROM category WHERE id = $1", categoryID)

	if err == sql.ErrNoRows {
		return model.Category{}, storage.ErrNotFound
	}

	if err != nil {
		return model.Category{}, fmt.Errorf("failed to get category: %w", err)
	}

	return c.toModel(), nil
}

func (p pg) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	if err := p.ext.GetContext(ctx, &category.ID, "INSERT INTO category (name) VALUES ($1) RETURNING id", category.Name); err != nil {
		return model.Category{}, fmt.Errorf("failed to create category: %w", err)
	}
	return category, nil
}

func (p pg) UpdateCategory(ctx context.Context, category model.Category) error {
	res, err := p.ext.ExecContext(ctx, "UPDATE category SET name = $1 WHERE id = $2", category.Name, category.ID)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) DeleteCategory(ctx context.Context, categoryID int64) error {
	res, err := p.ext.ExecContext(ctx, "DELETE FROM category WHERE id = $1", categoryID)

	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}
