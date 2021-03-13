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
	if err := p.ext.GetContext(ctx, &user.ID, `
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
	err := p.ext.GetContext(ctx, &u, "SELECT id, email, password_hash, is_admin FROM store_user WHERE email = $1", email)

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
	err := p.ext.GetContext(ctx, &u, "SELECT id, email, password_hash, is_admin FROM store_user WHERE id = $1", id)

	if err == sql.ErrNoRows {
		return model.User{}, storage.ErrNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return u.toModel(), nil
}

func (p pg) SaveToken(ctx context.Context, tokenID string, userID int64, expiresAt time.Time) error {
	if _, err := p.ext.ExecContext(ctx, `
			INSERT INTO token (id, user_id, expires_at) VALUES ($1, $2, $3)
		`, tokenID, userID, expiresAt); err != nil {

		return fmt.Errorf("failed to save token: %w", err)
	}
	return nil
}

func (p pg) DeleteToken(ctx context.Context, tokenID string) error {
	res, err := p.ext.ExecContext(ctx, "DELETE FROM token WHERE id = $1", tokenID)

	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) UpdateUserPermissions(ctx context.Context, user model.User) error {
	res, err := p.ext.ExecContext(ctx, `
		UPDATE store_user SET is_admin = $2 WHERE id = $1
	`, user.ID, user.IsAdmin)

	if err != nil {
		return fmt.Errorf("failed to update user permissions: %w", err)
	}

	if c, _ := res.RowsAffected(); c == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (p pg) InTx(ctx context.Context, action func(s storage.UserStorage) error) error {
	tx, err := p.dbx.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = action(pg{dbx: p.dbx, ext: tx})

	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %v root: %w", rbErr, err)
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
