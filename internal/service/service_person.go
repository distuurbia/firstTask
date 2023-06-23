// Package service realize bisnes-logic of the microservice
package service

import (
	"context"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
)

// Repository is an interface that contains CRUD methods and GetAll
type RepositoryPerson interface {
	Create(ctx context.Context, pers *model.Person) error
	ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Update(ctx context.Context, pers *model.Person) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ServicePerson contains Repository interface
type ServicePerson struct {
	persRps RepositoryPerson
}

// NewService accepts Repository object and returnes an object of type *PersonService
func NewServicePerson(persRps RepositoryPerson) *ServicePerson {
	return &ServicePerson{persRps: persRps}
}

// Create is a method of PersonService that calls Create method of Repository
func (srv *ServicePerson) Create(ctx context.Context, pers *model.Person) error {
	return srv.persRps.Create(ctx, pers)
}

// ReadRow is a method of PersonService that calls ReadRow method of Repository
func (srv *ServicePerson) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	return srv.persRps.ReadRow(ctx, id)
}

// Update is a method of PersonService that calls Update method of Repository
func (srv *ServicePerson) Update(ctx context.Context, pers *model.Person) error {
	return srv.persRps.Update(ctx, pers)
}

// Delete is a method of PersonService that calls Delete method of Repository
func (srv ServicePerson) Delete(ctx context.Context, id uuid.UUID) error {
	return srv.persRps.Delete(ctx, id)
}

// GetAll is a method of PersonService that calls  method of Repository
func (srv *ServicePerson) GetAll(ctx context.Context) ([]model.Person, error) {
	return srv.persRps.GetAll(ctx)
}