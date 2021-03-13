//+build integration

package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

func (s *postgresTestSuite) TestPg_CreateUser() {
	user := model.User{
		Email:        "admin@test.com",
		PasswordHash: "1234",
		IsAdmin:      true,
	}

	user, err := s.s.(pg).CreateUser(s.ctx, user)
	s.Require().NoError(err)

	s.Require().True(user.ID > 0, "ID is not populated")

	r := s.db.QueryRow("SELECT id, email, password_hash, is_admin FROM store_user WHERE id = $1", user.ID)
	res := model.User{}
	err = r.Scan(&res.ID, &res.Email, &res.PasswordHash, &res.IsAdmin)
	s.Require().NoError(err)

	s.Equal(user, res)

	_, err = s.s.(pg).CreateUser(s.ctx, model.User{Email: "admin@test.com", PasswordHash: "123"})
	s.True(errors.Is(storage.ErrEmailIsTaken, err))
}

func (s *postgresTestSuite) TestPg_GetUserByEmail() {
	_, err := s.db.Exec(`INSERT INTO store_user (email, password_hash, is_admin) VALUES ('admin@test.com', '123', TRUE);`)
	s.Require().NoError(err)

	u, err := s.s.(pg).GetUserByEmail(s.ctx, "admin@test.com")
	s.Require().NoError(err)

	s.Equal(model.User{ID: 1, Email: "admin@test.com", PasswordHash: "123", IsAdmin: true}, u)

	_, err = s.s.(pg).GetUserByEmail(s.ctx, "none@test.com")
	s.True(errors.Is(storage.ErrNotFound, err))
}

func (s *postgresTestSuite) TestPg_GetUserByID() {
	_, err := s.db.Exec(`INSERT INTO store_user (email, password_hash, is_admin) VALUES ('admin@test.com', '123', TRUE);`)
	s.Require().NoError(err)

	u, err := s.s.(pg).GetUserByID(s.ctx, 1)
	s.Require().NoError(err)

	s.Equal(model.User{ID: 1, Email: "admin@test.com", PasswordHash: "123", IsAdmin: true}, u)

	_, err = s.s.(pg).GetUserByID(s.ctx, 100500)
	s.True(errors.Is(storage.ErrNotFound, err))
}

func (s *postgresTestSuite) TestPg_SaveToken() {
	_, err := s.db.Exec(`INSERT INTO store_user (email, password_hash, is_admin) VALUES ('admin@test.com', '123', TRUE);`)
	s.Require().NoError(err)

	tokenID := uuid.NewString()
	userID := int64(1)
	expiredAt := time.Now().Add(10 * time.Hour).UTC()

	err = s.s.(pg).SaveToken(s.ctx, tokenID, userID, expiredAt)
	s.Require().NoError(err)

	r := s.db.QueryRow("SELECT user_id, expires_at FROM token WHERE id = $1", tokenID)
	var resUserID int64
	var resExpiresAt time.Time
	err = r.Scan(&resUserID, &resExpiresAt)
	s.Require().NoError(err)

	s.Equal(userID, resUserID)
	s.Equal(expiredAt, resExpiresAt.UTC())
}

func (s *postgresTestSuite) TestPg_DeleteToken() {
	_, err := s.db.Exec(`
		INSERT INTO store_user (email, password_hash, is_admin) VALUES ('admin@test.com', '123', TRUE);
		INSERT INTO token (id, user_id, expires_at) VALUES ('0e37df36-f698-11e6-8dd4-cb9ced3df976', 1, '2025-10-19 10:23:54')
	`)
	s.Require().NoError(err)

	err = s.s.(pg).DeleteToken(s.ctx, "0e37df36-f698-11e6-8dd4-cb9ced3df976")
	s.Require().NoError(err)

	err = s.s.(pg).DeleteToken(s.ctx, "0e37df36-f698-11e6-8dd4-cb9ced3df976")
	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_UpdateUserPermissions() {
	_, err := s.db.Exec(`
		INSERT INTO store_user (email, password_hash, is_admin) VALUES ('admin@test.com', '123', FALSE);
	`)
	s.Require().NoError(err)

	u := model.User{ID: 1, IsAdmin: true}

	err = s.s.(pg).UpdateUserPermissions(s.ctx, u)
	s.Require().NoError(err)

	res := model.User{}
	s.Require().NoError(s.db.QueryRow(`
		SELECT email, password_hash, is_admin FROM store_user WHERE id = $1
	`, 1).Scan(&res.Email, &res.PasswordHash, &res.IsAdmin))

	s.Equal(model.User{Email: "admin@test.com", PasswordHash: "123", IsAdmin: true}, res)

	err = s.s.(pg).UpdateUserPermissions(s.ctx, model.User{ID: 1000})
	s.True(errors.Is(err, storage.ErrNotFound))
}

func (s *postgresTestSuite) TestPg_InTx() {
	userCount := func() int {
		var c int
		s.Require().NoError(s.db.QueryRow(`SELECT count(*) FROM store_user`).Scan(&c))
		return c
	}

	baseline := userCount()

	err := s.s.(pg).InTx(s.ctx, func(us storage.UserStorage) error {
		_, err := us.CreateUser(s.ctx, model.User{Email: "user1@test.com", PasswordHash: "123"})
		s.Require().NoError(err)

		s.Equal(baseline, userCount(), "read uncommitted")

		return nil
	})

	s.Require().NoError(err)
	s.Equal(baseline+1, userCount(), "missing commit")

	baseline = userCount()

	err = s.s.(pg).InTx(s.ctx, func(us storage.UserStorage) error {
		_, err := us.CreateUser(s.ctx, model.User{Email: "user2@test.com", PasswordHash: "123"})
		s.Require().NoError(err)

		return assert.AnError
	})

	s.True(errors.Is(err, assert.AnError), fmt.Sprintf("wanted %s got %s", assert.AnError, err))
	s.Equal(baseline, userCount(), "missing rollback")
}
