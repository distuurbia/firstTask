package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonPgx struct {
	db *pgxpool.Pool
}

func NewPgxRep(db *pgxpool.Pool) *PersonPgx {
	return &PersonPgx{db: db}
}

func (pgxRps *PersonPgx) Create(ctx context.Context, pers *model.Person) error {
	if pers == nil {
		return ErrNil
	}
	_, err := pgxRps.db.Exec(ctx, "INSERT INTO persondb(salary, married, profession, id) VALUES($1, $2, $3, $4)", pers.Salary, pers.Married, pers.Profession, pers.Id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (pgxRps *PersonPgx) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	var pers model.Person
	err := pgxRps.db.QueryRow(ctx, "SELECT id, salary, married, profession FROM persondb WHERE id = $1", id).Scan(&pers.Id, &pers.Salary, &pers.Married, &pers.Profession)
	if err != nil {
		return &pers, fmt.Errorf("%w", err)
	}
	return &pers, nil
}

func (pgxRps *PersonPgx) Update(ctx context.Context, pers *model.Person) error {
	res, err := pgxRps.db.Exec(ctx, "UPDATE persondb SET salary = $1, married = $2, profession = $3 WHERE id = $4", pers.Salary, pers.Married, pers.Profession, pers.Id)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if res.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

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
