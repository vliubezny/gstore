package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (p pg) GetProducts(ctx context.Context, categoryID int64) ([]model.Product, error) {
	var products []product

	if err := p.db.SelectContext(ctx, &products, "SELECT * FROM product WHERE category_id=$1", categoryID); err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

	data := make([]model.Product, len(products))
	for i, d := range products {
		data[i] = d.toModel()
	}

	return data, nil
}

func (p pg) GetProduct(ctx context.Context, productID int64) (model.Product, error) {
	var prod product
	err := p.db.GetContext(ctx, &prod, "SELECT * FROM product WHERE id = $1", productID)

	if err == sql.ErrNoRows {
		return model.Product{}, storage.ErrNotFound
	}

	if err != nil {
		return model.Product{}, fmt.Errorf("failed to get product: %w", err)
	}

	return prod.toModel(), nil
}

func (p pg) CreateProduct(ctx context.Context, product model.Product) (model.Product, error) {
	if err := p.db.GetContext(ctx, &product.ID, `
			INSERT INTO product (category_id, name, description) VALUES ($1, $2, $3) RETURNING id
		`, product.CategoryID, product.Name, product.Description); err != nil {
		return model.Product{}, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

func (p pg) UpdateProduct(ctx context.Context, product model.Product) error {
	res, err := p.db.ExecContext(ctx, `
		UPDATE product SET
		category_id =$2,
		name = $3,
		description = $4
		WHERE id = $1
	`, product.ID, product.CategoryID, product.Name, product.Description)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) DeleteProduct(ctx context.Context, productID int64) error {
	res, err := p.db.ExecContext(ctx, "DELETE FROM product WHERE id = $1", productID)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}
