//+build integration

package postgres

import (
	"errors"

	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (s *postgresTestSuite) TestPg_GetStores() {
	_, err := s.db.Exec(`INSERT INTO store (name) VALUES
	('iStore'),
	('Amazon');`)
	s.Require().NoError(err)

	stores, err := s.s.GetStores(s.ctx)
	s.Require().NoError(err)

	s.Equal([]*model.Store{
		{ID: 1, Name: "iStore"},
		{ID: 2, Name: "Amazon"},
	}, stores)
}

func (s *postgresTestSuite) TestPg_GetStore() {
	_, err := s.db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	s.Require().NoError(err)

	store, err := s.s.GetStore(s.ctx, 1)
	s.Require().NoError(err)

	s.Equal(&model.Store{ID: 1, Name: "iStore"}, store)
}

func (s *postgresTestSuite) TestPg_GetStore_ErrNotFound() {
	_, err := s.s.GetStore(s.ctx, 100500)

	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_CreateStore() {
	str := &model.Store{
		Name: "test store",
	}

	err := s.s.CreateStore(s.ctx, str)
	s.Require().NoError(err)

	s.Require().True(str.ID > 0, "ID is not populated")

	r := s.db.QueryRow("SELECT name FROM store WHERE id = $1", str.ID)
	var name string
	err = r.Scan(&name)
	s.Require().NoError(err)

	s.Equal(str.Name, name)
}

func (s *postgresTestSuite) TestPg_UpdateStore() {
	_, err := s.db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	s.Require().NoError(err)

	str := &model.Store{
		ID:   1,
		Name: "test store",
	}

	err = s.s.UpdateStore(s.ctx, str)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT name FROM store WHERE id = $1", str.ID)
	var name string
	err = r.Scan(&name)
	s.Require().NoError(err)

	s.Equal(str.Name, name)
}

func (s *postgresTestSuite) TestPg_UpdateStore_ErrNotFound() {
	str := &model.Store{
		ID:   100500,
		Name: "test store",
	}

	err := s.s.UpdateStore(s.ctx, str)

	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_DeleteStore() {
	_, err := s.db.Exec(`INSERT INTO store (name) VALUES ('iStore');`)
	s.Require().NoError(err)
	id := int64(1)

	err = s.s.DeleteStore(s.ctx, id)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT count(*) FROM store WHERE id = $1", id)
	var c int
	err = r.Scan(&c)
	s.Require().NoError(err)

	s.Equal(0, c)
}

func (s *postgresTestSuite) TestPg_DeleteStore_ErrNotFound() {
	err := s.s.DeleteStore(s.ctx, 100500)

	s.True(errors.Is(err, storage.ErrNotFound))
}
