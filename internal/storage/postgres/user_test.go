//+build integration

package postgres

import (
	"errors"

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
