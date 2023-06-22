//Package repository is a package for work with db methods
package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
)

// SignIn create new user in users table
func (pgxRps *PersonPgx) SignIn(ctx context.Context, user *model.User) error {
	if user == nil {
		return ErrNil
	}
	var numberPersons int
	err := pgxRps.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE username = $1", user.Username).Scan(&numberPersons)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if numberPersons != 0 {
		return ErrExist
	}
	_, err = pgxRps.db.Exec(ctx, "INSERT INTO users(id, username, password) VALUES($1, $2, $3)", user.ID, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// SignIn verifies user with checking hash of the password from users table
func (pgxRps *PersonPgx) GetPassword(ctx context.Context, user *model.User) (string, error) {
	var hash string
	err := pgxRps.db.QueryRow(ctx, "SELECT password FROM users WHERE username = $1", user.Username).Scan(&hash)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return hash, nil
}

func (pgxRps *PersonPgx) AddRefreshToken(ctx context.Context, user *model.User) error{
	_, err := pgxRps.db.Exec(ctx, "UPDATE users SET refreshtoken = $1 WHERE username = $2", user.RefreshToken,  user.Username)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}