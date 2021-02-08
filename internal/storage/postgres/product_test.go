//+build integration

package postgres

import (
	"errors"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (s *postgresTestSuite) TestPg_GetProducts() {
	categoryID := int64(1)

	_, err := s.db.Exec(`INSERT INTO product (category_id, name, description) VALUES
	($1, 'iPhone 11', 'Old iphone'),
	($1, 'iPhone 12', 'New iphone');`, categoryID)
	s.Require().NoError(err)

	products, err := s.s.GetProducts(s.ctx, categoryID)
	s.Require().NoError(err)

	s.Equal([]model.Product{
		{ID: 1, CategoryID: categoryID, Name: "iPhone 11", Description: "Old iphone"},
		{ID: 2, CategoryID: categoryID, Name: "iPhone 12", Description: "New iphone"},
	}, products)
}

func (s *postgresTestSuite) TestPg_GetProduct() {
	_, err := s.db.Exec(`INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');`)
	s.Require().NoError(err)

	product, err := s.s.GetProduct(s.ctx, 1)
	s.Require().NoError(err)

	s.Equal(model.Product{ID: 1, CategoryID: 1, Name: "iPhone 11", Description: "Old iphone"}, product)
}

func (s *postgresTestSuite) TestPg_GetProduct_ErrNotFound() {
	_, err := s.s.GetProduct(s.ctx, 100500)

	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_CreateProduct() {
	prod := model.Product{
		CategoryID:  1,
		Name:        "iPhone 11",
		Description: "Old iphone",
	}

	prod, err := s.s.CreateProduct(s.ctx, prod)
	s.Require().NoError(err)

	s.Require().True(prod.ID > 0, "ID is not populated")

	r := s.db.QueryRow("SELECT id, category_id, name, description FROM product WHERE id = $1", prod.ID)
	res := model.Product{}
	err = r.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Description)
	s.Require().NoError(err)

	s.Equal(prod, res)

	_, err = s.s.CreateProduct(s.ctx, model.Product{CategoryID: 100, Name: "test", Description: "test"})
	s.True(errors.Is(storage.ErrUnknownCategory, err))
}

func (s *postgresTestSuite) TestPg_UpdateProduct() {
	_, err := s.db.Exec(`INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');`)
	s.Require().NoError(err)

	prod := model.Product{
		ID:          1,
		CategoryID:  2,
		Name:        "iPhone 12",
		Description: "New iphone",
	}

	err = s.s.UpdateProduct(s.ctx, prod)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT id, category_id, name, description FROM product WHERE id = $1", prod.ID)
	res := model.Product{}
	err = r.Scan(&res.ID, &res.CategoryID, &res.Name, &res.Description)
	s.Require().NoError(err)

	s.Equal(prod, res)

	s.True(errors.Is(storage.ErrNotFound, s.s.UpdateProduct(s.ctx,
		model.Product{ID: 100, CategoryID: 2, Name: "test", Description: "test"},
	)))

	s.True(errors.Is(storage.ErrUnknownCategory, s.s.UpdateProduct(s.ctx,
		model.Product{ID: 1, CategoryID: 100, Name: "test", Description: "test"},
	)))
}

func (s *postgresTestSuite) TestPg_DeleteProduct() {
	_, err := s.db.Exec(`INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');`)
	s.Require().NoError(err)

	id := int64(1)

	err = s.s.DeleteProduct(s.ctx, id)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT count(*) FROM product WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	s.Require().NoError(err)

	s.Equal(0, c)
}

func (s *postgresTestSuite) TestPg_DeleteProduct_ErrNotFound() {
	err := s.s.DeleteProduct(s.ctx, 100500)

	s.True(errors.Is(err, storage.ErrNotFound))
}
