// Package repository is a package for work with db methods
package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PersonPgx contains an object of type *pgxpool.Pool
type PersonPgx struct {
	db *pgxpool.Pool
}

// NewPgxRep accepts an object of *pgxpool.Pool and returns an object of type *PersonPgx
func NewPgxRep(db *pgxpool.Pool) *PersonPgx {
	return &PersonPgx{db: db}
}

// Create creates a row in postgreSQL
func (pgxRps *PersonPgx) Create(ctx context.Context, pers *model.Person) error {
	if pers == nil {
		return ErrNil
	}
	_, err := pgxRps.db.Exec(ctx, "INSERT INTO persondb(salary, married, profession, id) VALUES($1, $2, $3, $4)", pers.Salary, pers.Married, pers.Profession, pers.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

// ReadRow reads a row from postgreSQL
func (pgxRps *PersonPgx) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	var pers model.Person
	err := pgxRps.db.QueryRow(ctx, "SELECT id, salary, married, profession FROM persondb WHERE id = $1", id).Scan(&pers.ID, &pers.Salary, &pers.Married, &pers.Profession)
	if err != nil {
		return &pers, fmt.Errorf("%w", err)
	}
	return &pers, nil
}

// GetAll reads an all rows in postgreSQL
func (pgxRps *PersonPgx) GetAll(ctx context.Context) ([]model.Person, error) {
	var allPers []model.Person
	rows, err := pgxRps.db.Query(ctx, "SELECT id, salary, married, profession FROM persondb")
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer rows.Close()
	var pers model.Person
	for rows.Next() {
		err = rows.Scan(&pers.ID, &pers.Salary, &pers.Married, &pers.Profession)
		if err != nil {
			return allPers, fmt.Errorf("%w", err)
		}
		allPers = append(allPers, pers)
	}
	return allPers, nil
}

// Update updates a row in postgreSQL
func (pgxRps *PersonPgx) Update(ctx context.Context, pers *model.Person) error {
	res, err := pgxRps.db.Exec(ctx, "UPDATE persondb SET salary = $1, married = $2, profession = $3 WHERE id = $4", pers.Salary, pers.Married, pers.Profession, pers.ID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Delete deletes a row in postgreSQL
func (pgxRps *PersonPgx) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := pgxRps.db.Exec(ctx, "DELETE FROM persondb WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
