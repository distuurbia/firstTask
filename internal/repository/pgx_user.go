// Package repository is a package for work with db methods
package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
)

// SignUp creates new user in users table
func (rpsPgx *Pgx) SignUp(ctx context.Context, user *model.User) error {
	if user == nil {
		return ErrNil
	}
	var numberPersons int
	err := rpsPgx.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = $1", user.Username).Scan(&numberPersons)
	if err != nil {
		return fmt.Errorf("Pgx -> SignUp -> QueryRow -> error: %w", err)
	}
	if numberPersons != 0 {
		return ErrExist
	}
	_, err = rpsPgx.db.Exec(ctx, "INSERT INTO users(id, username, password) VALUES($1, $2, $3)", user.ID, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("Pgx -> SignUp -> Exec -> error: %w", err)
	}
	return nil
}

// GetPasswordAndIDByUsername hash of the password and id from users table
func (rpsPgx *Pgx) GetPasswordAndIDByUsername(ctx context.Context, username string) (uuid.UUID, []byte, error) {
	var user model.User
	user.Username = username
	err := rpsPgx.db.QueryRow(ctx, "SELECT id, password FROM users WHERE username = $1", user.Username).Scan(&user.ID, &user.Password)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("Pgx -> GetPasswordAndIDByUserName -> QueryRow -> error: %w", err)
	}
	return user.ID, user.Password, nil
}

// GetRefreshTokenByID returnes refreshToken from users table by id
func (rpsPgx *Pgx) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error) {
	var hash string
	err := rpsPgx.db.QueryRow(ctx, "SELECT refreshToken FROM users WHERE id = $1", id).Scan(&hash)
	if err != nil {
		return "", fmt.Errorf("Pgx -> GetRefreshTokenByName -> QueryRow -> error: %w", err)
	}
	return hash, nil
}

// AddRefreshToken adds refreshToken to users table by id
func (rpsPgx *Pgx) AddRefreshToken(ctx context.Context, user *model.User) error {
	_, err := rpsPgx.db.Exec(ctx, "UPDATE users SET refreshtoken = $1 WHERE id = $2", user.RefreshToken, user.ID)
	if err != nil {
		return fmt.Errorf("Pgx -> AddRefreshToken -> Exec -> error: %w", err)
	}
	return nil
}
