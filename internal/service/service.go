package service

import (
	"context"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
)
type Repository interface{
	Create(ctx context.Context, pers *model.Person) error;
	ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error);
	Update(ctx context.Context, pers *model.Person) error;
	Delete(ctx context.Context, id uuid.UUID) error 
}
type PersonService struct {
	rps Repository
}

func NewService(rps Repository) *PersonService {

	return &PersonService{rps: rps}

}

func (srv *PersonService) Create(ctx context.Context, pers *model.Person) error {

	return srv.rps.Create(ctx, pers)

}
func (srv *PersonService) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {

	return srv.rps.ReadRow(ctx, id)

}
func (srv *PersonService) Update(ctx context.Context, pers *model.Person) error {

	return srv.rps.Update(ctx, pers)

}
func (srv *PersonService) Delete(ctx context.Context, id uuid.UUID) error {

	return srv.rps.Delete(ctx, id)

}
