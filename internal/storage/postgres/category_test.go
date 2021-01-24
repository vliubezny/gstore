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

func TestPg_GetCategories(t *testing.T) {
	categories, err := s.GetCategories(ctx)
	require.NoError(t, err)

	assert.Equal(t, []model.Category{
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

func TestPg_GetCategory(t *testing.T) {
	category, err := s.GetCategory(ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, model.Category{ID: 1, Name: "Electronics"}, category)
}

func TestPg_GetCategory_ErrNotFound(t *testing.T) {
	_, err := s.GetCategory(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_CreateCategory(t *testing.T) {
	c := model.Category{
		Name: "test category",
	}

	c, err := s.CreateCategory(ctx, c)
	require.NoError(t, err)

	require.True(t, c.ID > 0, "ID is not populated")

	r := db.QueryRow("SELECT name FROM category WHERE id = $1", c.ID)
	var name string
	err = r.Scan(&name)
	require.NoError(t, err)

	assert.Equal(t, c.Name, name)
}

func TestPg_UpdateCategory(t *testing.T) {
	c := model.Category{
		ID:   1,
		Name: "test category",
	}

	err := s.UpdateCategory(ctx, c)
	require.NoError(t, err)

	r := db.QueryRow("SELECT name FROM category WHERE id = $1", c.ID)
	var name string
	err = r.Scan(&name)
	require.NoError(t, err)

	assert.Equal(t, c.Name, name)
}

func TestPg_UpdateCategory_ErrNotFound(t *testing.T) {
	c := model.Category{
		ID:   100500,
		Name: "test category",
	}

	err := s.UpdateCategory(ctx, c)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_DeleteCategory(t *testing.T) {
	var id int64 = 1

	err := s.DeleteCategory(ctx, id)
	require.NoError(t, err)

	r := db.QueryRow("SELECT count(*) FROM category WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	require.NoError(t, err)

	assert.Equal(t, 0, c)
}

func TestPg_DeleteCategory_ErrNotFound(t *testing.T) {
	err := s.DeleteCategory(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}
