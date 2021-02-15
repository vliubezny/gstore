package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

const emailUniqueConstraint = "store_user_email_key"

func (p pg) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	if err := p.db.GetContext(ctx, &user.ID, `
			INSERT INTO store_user (email, password_hash, is_admin) VALUES ($1, $2, $3) RETURNING id
		`, user.Email, user.PasswordHash, user.IsAdmin); err != nil {

		if err, ok := err.(*pq.Error); ok && err.Constraint == emailUniqueConstraint {
			return model.User{}, storage.ErrEmailIsTaken
		}
		return model.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (p pg) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	var u user
	err := p.db.GetContext(ctx, &u, "SELECT * FROM store_user WHERE email = $1", email)

	if err == sql.ErrNoRows {
		return model.User{}, storage.ErrNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toModel(), nil
}
