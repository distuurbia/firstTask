package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Person struct {
	db *pgxpool.Pool
}

func NewPerson(db *pgxpool.Pool) *Person {
	return &Person{db: db}
}

func (r *Person) Create(ctx context.Context) error {

}

func (r *Person) Read() error {}

func (r *Person) Update() error {}

func (r *Person) Delete() error {}
