package service

import (
	"context"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/repository"
	"github.com/google/uuid"
)

type PersonService struct {
	persPgx *repository.PersonPgx
}

func NewService(persPgx *repository.PersonPgx) *PersonService {

	return &PersonService{persPgx: persPgx}

}

func (srv *PersonService) Create(ctx context.Context, p *model.Person) error {

	return srv.persPgx.Create(ctx, p)

}
func (srv *PersonService) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {

	return srv.persPgx.ReadRow(ctx, id)

}
func (srv *PersonService) Update(ctx context.Context, p *model.Person) error {

	return srv.persPgx.Update(ctx, p)

}
func (srv *PersonService) Delete(ctx context.Context, id uuid.UUID) error {

	return srv.persPgx.Delete(ctx, id)

}
