//+build integration

package postgres

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (s *postgresTestSuite) TestPg_GetStorePositions() {
	_, err := s.db.Exec(`
		INSERT INTO product (category_id, name, description) VALUES
			(1, 'iPhone 11', 'Old iphone'),
			(1, 'iPhone 12', 'New iphone');
		INSERT INTO store (name) VALUES ('iStore');
		INSERT INTO position (product_id, store_id, price) VALUES
			(1, 1, 100),
			(2, 1, 200);
	`)
	s.Require().NoError(err)
	storeID := int64(1)

	positions, err := s.s.GetStorePositions(s.ctx, storeID)
	s.Require().NoError(err)

	s.Equal([]model.Position{
		{ProductID: 1, StoreID: storeID, Price: decimal.NewFromInt(100)},
		{ProductID: 2, StoreID: storeID, Price: decimal.NewFromInt(200)},
	}, positions)
}

func (s *postgresTestSuite) TestPg_GetProductPositions() {
	_, err := s.db.Exec(`
		INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');
		INSERT INTO store (name) VALUES ('iStore'), ('Amazon');
		INSERT INTO position (product_id, store_id, price) VALUES
			(1, 1, 100),
			(1, 2, 200);
	`)
	s.Require().NoError(err)
	productID := int64(1)

	positions, err := s.s.GetProductPositions(s.ctx, productID)
	s.Require().NoError(err)

	s.Equal([]model.Position{
		{ProductID: productID, StoreID: 1, Price: decimal.NewFromInt(100)},
		{ProductID: productID, StoreID: 2, Price: decimal.NewFromInt(200)},
	}, positions)
}

func (s *postgresTestSuite) TestPg_UpsertPosition() {
	_, err := s.db.Exec(`
		INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');
		INSERT INTO store (name) VALUES ('iStore');
	`)
	s.Require().NoError(err)

	p := model.Position{
		ProductID: 1,
		StoreID:   1,
		Price:     decimal.NewFromInt(100),
	}

	s.Require().NoError(s.s.UpsertPosition(s.ctx, p))

	r := s.db.QueryRow(`
		SELECT product_id, store_id, price FROM position WHERE product_id = $1 AND store_id = $2
	`, p.ProductID, p.StoreID)
	res := model.Position{}
	err = r.Scan(&res.ProductID, &res.StoreID, &res.Price)
	s.Require().NoError(err)

	s.Equal(p, res)

	p.Price = decimal.NewFromInt(200)

	s.Require().NoError(s.s.UpsertPosition(s.ctx, p))

	r = s.db.QueryRow(`
		SELECT product_id, store_id, price FROM position WHERE product_id = $1 AND store_id = $2;
	`, p.ProductID, p.StoreID)
	err = r.Scan(&res.ProductID, &res.StoreID, &res.Price)
	s.Require().NoError(err)

	s.Equal(p, res)

	p.Price = decimal.Zero

	s.True(errors.Is(storage.ErrInvalidPrice, s.s.UpsertPosition(s.ctx, p)))
}

func (s *postgresTestSuite) TestPg_DeletePosition() {
	_, err := s.db.Exec(`
		INSERT INTO product (category_id, name, description) VALUES (1, 'iPhone 11', 'Old iphone');
		INSERT INTO store (name) VALUES ('iStore');
		INSERT INTO position (product_id, store_id, price) VALUES (1, 1, 100);
	`)
	s.Require().NoError(err)

	productID := int64(1)
	storeID := int64(1)

	err = s.s.DeletePosition(s.ctx, productID, storeID)
	s.Require().NoError(err)

	r := s.db.QueryRow(`
		SELECT count(*) FROM position WHERE product_id = $1 AND store_id = $2;
	`, productID, storeID)
	var c int
	err = r.Scan(&c)
	s.Require().NoError(err)

	s.Equal(0, c)

	err = s.s.DeletePosition(s.ctx, 100, 500)

	s.True(errors.Is(err, storage.ErrNotFound))
}
