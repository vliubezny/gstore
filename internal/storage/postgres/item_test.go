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

func TestPg_GetStoreItems(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;
			ALTER SEQUENCE item_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	var storeID int64
	err := db.QueryRow(`INSERT INTO store (name) VALUES ('iStore') RETURNING id;`).Scan(&storeID)
	require.NoError(t, err)

	_, err = db.Exec(`INSERT INTO item (store_id, name, description, price) VALUES
	($1, 'iPhone 11', 'Old iphone', 100000),
	($1, 'iPhone 12', 'New iphone', 200000);`, storeID)
	require.NoError(t, err)

	items, err := s.GetStoreItems(ctx, storeID)
	require.NoError(t, err)

	assert.Equal(t, []*model.Item{
		{ID: 1, StoreID: storeID, Name: "iPhone 11", Description: "Old iphone", Price: 100000},
		{ID: 2, StoreID: storeID, Name: "iPhone 12", Description: "New iphone", Price: 200000},
	}, items)
}

func TestPg_GetStoreItem(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;
			ALTER SEQUENCE item_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore');
		INSERT INTO item (store_id, name, description, price) VALUES (1, 'iPhone 11', 'Old iphone', 100000);`)
	require.NoError(t, err)

	item, err := s.GetStoreItem(ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, &model.Item{ID: 1, StoreID: 1, Name: "iPhone 11", Description: "Old iphone", Price: 100000}, item)
}

func TestPg_GetStoreItem_ErrNotFound(t *testing.T) {
	_, err := s.GetStoreItem(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_CreateStoreItem(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;
			ALTER SEQUENCE item_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	require.NoError(t, err)

	itm := &model.Item{
		ID:          1,
		StoreID:     1,
		Name:        "iPhone 11",
		Description: "Old iphone",
		Price:       100000,
	}

	err = s.CreateStoreItem(ctx, itm)
	require.NoError(t, err)

	require.True(t, itm.ID > 0, "ID is not populated")

	r := db.QueryRow("SELECT id, store_id, name, description, price FROM item WHERE id = $1", itm.ID)
	res := &model.Item{}
	err = r.Scan(&res.ID, &res.StoreID, &res.Name, &res.Description, &res.Price)
	require.NoError(t, err)

	assert.Equal(t, itm, res)
}

func TestPg_UpdateStoreItem(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;
			ALTER SEQUENCE item_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore'), ('Amazon');
		INSERT INTO item (store_id, name, description, price) VALUES (1, 'iPhone 11', 'Old iphone', 100000);`)
	require.NoError(t, err)

	itm := &model.Item{
		ID:          1,
		StoreID:     2,
		Name:        "iPhone 12",
		Description: "New iphone",
		Price:       200000,
	}

	err = s.UpdateStoreItem(ctx, itm)
	require.NoError(t, err)

	r := db.QueryRow("SELECT id, store_id, name, description, price FROM item WHERE id = $1", itm.ID)
	res := &model.Item{}
	err = r.Scan(&res.ID, &res.StoreID, &res.Name, &res.Description, &res.Price)
	require.NoError(t, err)

	assert.Equal(t, itm, res)
}

func TestPg_UpdateStoreItem_ErrNotFound(t *testing.T) {
	itm := &model.Item{
		ID:          100500,
		StoreID:     2,
		Name:        "iPhone 12",
		Description: "New iphone",
		Price:       200000,
	}

	err := s.UpdateStoreItem(ctx, itm)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_DeleteStoreItem(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;
			ALTER SEQUENCE item_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore'), ('Amazon');
		INSERT INTO item (store_id, name, description, price) VALUES (1, 'iPhone 11', 'Old iphone', 100000);`)
	require.NoError(t, err)

	id := int64(1)

	err = s.DeleteStoreItem(ctx, id)
	require.NoError(t, err)

	r := db.QueryRow("SELECT count(*) FROM item WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	require.NoError(t, err)

	assert.Equal(t, 0, c)
}

func TestPg_DeleteStoreItem_ErrNotFound(t *testing.T) {
	err := s.DeleteStoreItem(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}
