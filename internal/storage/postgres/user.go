package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (p pg) GetUserByID(ctx context.Context, id int64) (model.User, error) {
	var u user
	err := p.db.GetContext(ctx, &u, "SELECT * FROM store_user WHERE id = $1", id)

	if err == sql.ErrNoRows {
		return model.User{}, storage.ErrNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toModel(), nil
}

func (p pg) SaveToken(ctx context.Context, tokenID string, userID int64, expiresAt time.Time) error {
	if _, err := p.db.ExecContext(ctx, `
			INSERT INTO token (id, user_id, expires_at) VALUES ($1, $2, $3)
		`, tokenID, userID, expiresAt); err != nil {

		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

func (p pg) DeleteToken(ctx context.Context, tokenID string) error {
	res, err := p.db.ExecContext(ctx, "DELETE FROM token WHERE id = $1", tokenID)

	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}
