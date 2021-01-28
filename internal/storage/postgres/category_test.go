//+build integration

package postgres

import (
	"errors"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (s *postgresTestSuite) TestPg_GetCategories() {
	categories, err := s.s.GetCategories(s.ctx)
	s.Require().NoError(err)

	s.Equal([]model.Category{
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
	}, categories)
}

func (s *postgresTestSuite) TestPg_GetCategory() {
	category, err := s.s.GetCategory(s.ctx, 1)
	s.Require().NoError(err)

	s.Equal(model.Category{ID: 1, Name: "Electronics"}, category)
}

func (s *postgresTestSuite) TestPg_GetCategory_ErrNotFound() {
	_, err := s.s.GetCategory(s.ctx, 100500)

	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_CreateCategory() {
	c := model.Category{
		Name: "test category",
	}

	c, err := s.s.CreateCategory(s.ctx, c)
	s.Require().NoError(err)

	s.Require().True(c.ID > 0, "ID is not populated")

	r := s.db.QueryRow("SELECT name FROM category WHERE id = $1", c.ID)
	var name string
	err = r.Scan(&name)
	s.Require().NoError(err)

	s.Equal(c.Name, name)
}

func (s *postgresTestSuite) TestPg_UpdateCategory() {
	c := model.Category{
		ID:   1,
		Name: "test category",
	}

	err := s.s.UpdateCategory(s.ctx, c)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT name FROM category WHERE id = $1", c.ID)
	var name string
	err = r.Scan(&name)
	s.Require().NoError(err)

	s.Equal(c.Name, name)
}

func (s *postgresTestSuite) TestPg_UpdateCategory_ErrNotFound() {
	c := model.Category{
		ID:   100500,
		Name: "test category",
	}

	err := s.s.UpdateCategory(s.ctx, c)

	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_DeleteCategory() {
	var id int64 = 5

	err := s.s.DeleteCategory(s.ctx, id)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT count(*) FROM category WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	s.Require().NoError(err)

	s.Equal(0, c)
}

func (s *postgresTestSuite) TestPg_DeleteCategory_ErrNotFound() {
	err := s.s.DeleteCategory(s.ctx, 100500)

	s.True(errors.Is(err, storage.ErrNotFound))
}
