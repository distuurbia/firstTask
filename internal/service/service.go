package service

import (
	"context"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/repository/postgreRep"
	"github.com/google/uuid"
)

type Person struct {
	rps *repository.Person
}

func NewService(rps *repository.Person) *Person {

	return &Person{rps: rps}

}

func (s *Person) Create(ctx context.Context, p *model.Person) error {

	return s.rps.Create(ctx, p)

}
func (s *Person) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {

	return s.rps.ReadRow(ctx, id)

}
func (s *Person) Update(ctx context.Context, p *model.Person) error {

	return s.rps.Update(ctx, p)

}
func (s *Person) Delete(ctx context.Context, id uuid.UUID) error {

	return s.rps.Delete(ctx, id)

}
