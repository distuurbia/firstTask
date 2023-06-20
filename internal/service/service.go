// Package service realize bisnes-logic of the microservice
package service

import (
	"context"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
)

// Repository is an interface that contains CRUD methods and GetAll
type Repository interface {
	Create(ctx context.Context, pers *model.Person) error
	ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Update(ctx context.Context, pers *model.Person) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PersonService contains Repository interface
type PersonService struct {
	rps Repository
}

// NewService accepts Repository object and returnes an object of type *PersonService
func NewService(rps Repository) *PersonService {
	return &PersonService{rps: rps}
}

// Create is a method of PersonService that calls Create method of Repository
func (srv *PersonService) Create(ctx context.Context, pers *model.Person) error {
	return srv.rps.Create(ctx, pers)
}

// ReadRow is a method of PersonService that calls ReadRow method of Repository
func (srv *PersonService) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	return srv.rps.ReadRow(ctx, id)
}

// Update is a method of PersonService that calls Update method of Repository
func (srv *PersonService) Update(ctx context.Context, pers *model.Person) error {
	return srv.rps.Update(ctx, pers)
}

// Delete is a method of PersonService that calls Delete method of Repository
func (srv *PersonService) Delete(ctx context.Context, id uuid.UUID) error {
	return srv.rps.Delete(ctx, id)
}

// GetAll is a method of PersonService that calls  method of Repository
func (srv *PersonService) GetAll(ctx context.Context) ([]model.Person, error) {
	return srv.rps.GetAll(ctx)
}
