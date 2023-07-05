// Package service realize bisnes-logic of the microservice
package service

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// PersonRepository is an interface that contains CRUD methods and GetAll
type PersonRepository interface {
	Create(ctx context.Context, pers *model.Person) error
	ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error)
	GetAll(ctx context.Context) ([]model.Person, error)
	Update(ctx context.Context, pers *model.Person) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// PersonRedisRepository is an interface that contains redis methods
type PersonRedisRepository interface {
	Set(ctx context.Context, user *model.Person) error
	Get(ctx context.Context, id uuid.UUID) (*model.Person, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// AddToStream(ctx context.Context, pers *model.Person) error
	// GetFromStream(ctx context.Context, id uuid.UUID) (*model.Person, error)
	// DeleteFromStream(ctx context.Context, id uuid.UUID) error
	// UpdateToStream(ctx context.Context, pers *model.Person) error
}

// PersonService contains Repository interface
type PersonService struct {
	persRps    PersonRepository
	persRdsRps PersonRedisRepository
}

// NewPersonService accepts Repository object and returnes an object of type *PersonService
func NewPersonService(persRps PersonRepository, persRdsRps PersonRedisRepository) *PersonService {
	return &PersonService{persRps: persRps, persRdsRps: persRdsRps}
}

// Create is a method of PersonService that calls Create method of Repository
func (srv *PersonService) Create(ctx context.Context, pers *model.Person) error {
	err := srv.persRps.Create(ctx, pers)
	if err != nil {
		return fmt.Errorf("PersonService -> Create -> persRps.Create -> error: %w", err)
	}
	err = srv.persRdsRps.Set(ctx, pers)
	if err != nil {
		return fmt.Errorf("PersonService -> Create -> persRdsRps.Set -> error: %w", err)
	}
	return nil
}

// ReadRow is a method of PersonService that calls ReadRow method of Repository
func (srv *PersonService) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	pers, err := srv.persRdsRps.Get(ctx, id)
	if err != nil && err.Error() != redis.Nil.Error() {
		return nil, fmt.Errorf("PersonService -> ReadRow -> persRdsRps.Get -> error: %w", err)
	}
	if pers == nil {
		pers, err = srv.persRps.ReadRow(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("PersonService -> ReadRow -> persRps.ReadRow -> error: %w", err)
		}
		err = srv.persRdsRps.Set(ctx, pers)
		if err != nil && err != redis.Nil {
			return nil, fmt.Errorf("PersonService -> ReadRow -> persRdsRps.Set -> error: %w", err)
		}
	}
	return pers, nil
}

// Update is a method of PersonService that calls Update method of Repository
func (srv *PersonService) Update(ctx context.Context, pers *model.Person) error {
	err := srv.persRps.Update(ctx, pers)
	if err != nil {
		return fmt.Errorf("PersonService -> Update -> persRps -> Update -> error: %w", err)
	}
	err = srv.persRdsRps.Delete(ctx, pers.ID)
	if err != nil && err.Error() != redis.Nil.Error() {
		return fmt.Errorf("PersonService -> Update -> persRdsRps.Delete -> error: %w", err)
	}
	err = srv.persRdsRps.Set(ctx, pers)
	if err != nil {
		return fmt.Errorf("PersonService -> Update -> persRdsRps.Set -> error: %w", err)
	}
	return nil
}

// Delete is a method of PersonService that calls Delete method of Repository
func (srv PersonService) Delete(ctx context.Context, id uuid.UUID) error {
	err := srv.persRps.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("PersonService -> Delete -> persRps.Delete -> error: %w", err)
	}
	_ = srv.persRdsRps.Delete(ctx, id)
	if err != nil && err.Error() != redis.Nil.Error() {
		return fmt.Errorf("PersonService -> Delete -> persRdsRps.Delete -> error: %w", err)
	}
	return nil
}

// GetAll is a method of PersonService that calls  method of Repository
func (srv *PersonService) GetAll(ctx context.Context) ([]model.Person, error) {
	return srv.persRps.GetAll(ctx)
}
