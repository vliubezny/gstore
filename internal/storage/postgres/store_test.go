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

func TestPg_GetStores(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES
	('iStore'),
	('Amazon');`)
	require.NoError(t, err)

	stores, err := s.GetStores(ctx)
	require.NoError(t, err)

	assert.Equal(t, []*model.Store{
		{ID: 1, Name: "iStore"},
		{ID: 2, Name: "Amazon"},
	}, stores)
}

func TestPg_GetStore(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	require.NoError(t, err)

	store, err := s.GetStore(ctx, 1)
	require.NoError(t, err)

	assert.Equal(t, &model.Store{ID: 1, Name: "iStore"}, store)
}

func TestPg_GetStore_ErrNotFound(t *testing.T) {
	_, err := s.GetStore(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_CreateStore(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	str := &model.Store{
		Name: "test store",
	}

	err := s.CreateStore(ctx, str)
	require.NoError(t, err)

	require.True(t, str.ID > 0, "ID is not populated")

	r := db.QueryRow("SELECT name FROM store WHERE id = $1", str.ID)
	var name string
	err = r.Scan(&name)
	require.NoError(t, err)

	assert.Equal(t, str.Name, name)
}

func TestPg_UpdateStore(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	require.NoError(t, err)

	str := &model.Store{
		ID:   1,
		Name: "test store",
	}

	err = s.UpdateStore(ctx, str)
	require.NoError(t, err)

	r := db.QueryRow("SELECT name FROM store WHERE id = $1", str.ID)
	var name string
	err = r.Scan(&name)
	require.NoError(t, err)

	assert.Equal(t, str.Name, name)
}

func TestPg_UpdateStore_ErrNotFound(t *testing.T) {
	str := &model.Store{
		ID:   100500,
		Name: "test store",
	}

	err := s.UpdateStore(ctx, str)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}

func TestPg_DeleteStore(t *testing.T) {
	defer func() {
		_, err := db.Exec(`DELETE FROM store;
			ALTER SEQUENCE store_id_seq RESTART WITH 1;`)
		require.NoError(t, err)
	}()

	_, err := db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	require.NoError(t, err)
	id := int64(1)

	err = s.DeleteStore(ctx, id)
	require.NoError(t, err)

	r := db.QueryRow("SELECT count(*) FROM store WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	require.NoError(t, err)

	assert.Equal(t, 0, c)
}

func TestPg_DeleteStore_ErrNotFound(t *testing.T) {
	err := s.DeleteStore(ctx, 100500)

	assert.True(t, errors.Is(err, storage.ErrNotFound))
}
