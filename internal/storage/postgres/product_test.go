//+build integration

package postgres

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func TestPg_GetProducts(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM product;
		ALTER SEQUENCE product_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	categoryID := int64(1)

	_, err := db.Exec(`INSERT INTO product (category_id, name, description) VALUES
	($1, 'iPhone 11', 'Old iphone'),
	($1, 'iPhone 12', 'New iphone');`, categoryID)
	require.NoError(t, err)

	products, err := s.GetProducts(ctx, categoryID)
	require.NoError(t, err)

	assert.Equal(t, []*model.Product{
		{ID: 1, CategoryID: categoryID, Name: "iPhone 11", Description: "Old iphone"},
		{ID: 2, CategoryID: categoryID, Name: "iPhone 12", Description: "New iphone"},
	}, products)
}

func TestPg_GetProduct(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM product;
		ALTER SEQUENCE product_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');`)
	require.NoError(t, err)

	product, err := s.GetProduct(ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, &model.Product{ID: 1, CategoryID: 1, Name: "iPhone 11", Description: "Old iphone"}, product)
}

func TestPg_GetProduct_ErrNotFound(t *testing.T) {
	_, err := s.GetProduct(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_CreateProduct(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM product;
		ALTER SEQUENCE product_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	prod := &model.Product{
		ID:          1,
		CategoryID:  1,
		Name:        "iPhone 11",
		Description: "Old iphone",
	}

	err := s.CreateProduct(ctx, prod)
	require.NoError(t, err)

	require.True(t, prod.ID > 0, "ID is not populated")

	r := db.QueryRow("SELECT id, category_id, name, description FROM product WHERE id = $1", prod.ID)
	res := &model.Product{}
	err = r.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Description)
	require.NoError(t, err)

	assert.Equal(t, prod, res)
}

func TestPg_UpdateProduct(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM product;
		ALTER SEQUENCE product_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');`)
	require.NoError(t, err)

	prod := &model.Product{
		ID:          1,
		CategoryID:  2,
		Name:        "iPhone 12",
		Description: "New iphone",
	}

	err = s.UpdateProduct(ctx, prod)
	require.NoError(t, err)

	r := db.QueryRow("SELECT id, category_id, name, description FROM product WHERE id = $1", prod.ID)
	res := &model.Product{}
	err = r.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Description)
	require.NoError(t, err)

	assert.Equal(t, prod, res)
}

func TestPg_UpdateProduct_ErrNotFound(t *testing.T) {
	prod := &model.Product{
		ID:          100500,
		CategoryID:  2,
		Name:        "iPhone 12",
		Description: "New iphone",
	}

	err := s.UpdateProduct(ctx, prod)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_DeleteProduct(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM product;
		ALTER SEQUENCE product_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');`)
	require.NoError(t, err)

	id := int64(1)

	err = s.DeleteProduct(ctx, id)
	require.NoError(t, err)

	r := db.QueryRow("SELECT count(*) FROM product WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	require.NoError(t, err)

	assert.Equal(t, 0, c)
}

func TestPg_DeleteProduct_ErrNotFound(t *testing.T) {
	err := s.DeleteProduct(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}
