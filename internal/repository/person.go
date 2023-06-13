package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct {
	db *pgxpool.Pool
}

func NewPerson(db *pgxpool.Pool) *Person {
	return &Person{db: db}
}

func (r *Person) Create(ctx context.Context, p *model.Person) error {
	var id uuid.UUID
	id = uuid.New()
	_, err := r.db.Exec(ctx, "INSERT INTO persondb(salary, married, profession, id) VALUES($1, $2, $3, $4)", p.Salary, p.Married, p.Profession, id)
	if err != nil {
		return fmt.Errorf("create %w")
	}
	return nil
}

func (r *Person) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	var p model.Person
	err := r.db.QueryRow(ctx, "SELECT salary, married, profession FROM persondb WHERE id = $1", p.Id).Scan(&p)
	if err != nil {
		return &p, fmt.Errorf("readRow %w")
	}
	return &p, nil
}

func (r *Person) Update(ctx context.Context, p *model.Person, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "UPDATE persondb SET salary = $1, married = $2, profession = $3 WHERE id = $4", p.Salary, p.Married, p.Profession, id)
	if err != nil {
		return fmt.Errorf("Update %w")
	}
	return nil
}

func (r *Person) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("Delete %w")
	}
	return nil
}
