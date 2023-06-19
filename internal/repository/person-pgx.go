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

func NewRepository(db *pgxpool.Pool) *PersonPgx {
	return &PersonPgx{db: db}
}

func (r *PersonPgx) Create(ctx context.Context, p *model.Person) error {
	if p == nil {
		return ErrNil
	}
	_, err := r.db.Exec(ctx, "INSERT INTO persondb(salary, married, profession, id) VALUES($1, $2, $3, $4)", p.Salary, p.Married, p.Profession, p.Id)
	if err != nil {
		return fmt.Errorf("create %w", err)
	}
	return nil
}

func (r *PersonPgx) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	var p model.Person
	err := r.db.QueryRow(ctx, "SELECT id, salary, married, profession FROM persondb WHERE id = $1", id).Scan(&p.Id, &p.Salary, &p.Married, &p.Profession)
	if err != nil {
		return &p, fmt.Errorf("readRow %w", err)
	}
	return &p, nil
}

func (r *PersonPgx) Update(ctx context.Context, p *model.Person) error {
	_, err := r.ReadRow(ctx, p.Id)
	if err != nil{
		return pgx.ErrNoRows
	}
	_, err = r.db.Exec(ctx, "UPDATE persondb SET salary = $1, married = $2, profession = $3 WHERE id = $4", p.Salary, p.Married, p.Profession, p.Id)
	if err != nil {
		return fmt.Errorf("update %w", err)
	}
	return nil
}

func (r *PersonPgx) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.ReadRow(ctx, id)
	if err != nil{
		return pgx.ErrNoRows
	}
	_, err = r.db.Exec(ctx, "DELETE FROM persondb WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete %w", err)
	}
	return nil
}
