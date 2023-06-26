// handler contains handler methods and handler tests
package handler

import (
	"context"
	"os"
	"testing"

	"github.com/distuurbia/firstTask/internal/handler/mocks"
	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// vladimir is test object of model.Person
var vladimir = model.Person{
	ID:         uuid.New(),
	Salary:     300,
	Married:    false,
	Profession: "teacher",
}

// srvc is object of *mocks.Service
var srvc *mocks.Service

// TestMain execute all tests
func TestMain(m *testing.M) {
	srvc = new(mocks.Service)
	exitVal := m.Run()
	os.Exit(exitVal)
}

// TestCreate is a mocktest for Create method of interface Service
func TestCreate(t *testing.T) {
	srvc.On("Create", mock.Anything, mock.AnythingOfType("*model.Person")).Return(nil).Once()

	handl := service.NewPersonService(srvc)
	err := handl.Create(context.Background(), &vladimir)
	assert.NoError(t, err)

	srvc.AssertExpectations(t)
}

// TestReadRow is a mocktest for ReadRow method of interface Service
func TestReadRow(t *testing.T) {
	srvc.On("ReadRow", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&vladimir, nil)
	handle := service.NewPersonService(srvc)
	testVladimir, err := handle.ReadRow(context.Background(), vladimir.ID)
	assert.Nil(t, err)
	assert.Equal(t, testVladimir.ID, vladimir.ID)
	assert.Equal(t, testVladimir.Salary, vladimir.Salary)
	assert.Equal(t, testVladimir.Married, vladimir.Married)
	assert.Equal(t, testVladimir.Profession, vladimir.Profession)
	srvc.AssertExpectations(t)
}

func TestGetAll(t *testing.T) {
	srvc.On("GetAll", mock.Anything).Return([]model.Person{vladimir}, nil)
	handle := service.NewPersonService(srvc)
	allPers, err := handle.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, len(allPers), len([]model.Person{vladimir}))
}

// TestUpdate is a mocktest for Update method of interface Service
func TestUpdate(t *testing.T) {
	srvc.On("Update", mock.Anything, mock.AnythingOfType("*model.Person")).Return(nil).Once()
	vladimir.Salary = 700
	vladimir.Married = false
	vladimir.Profession = "Lawer"
	handle := service.NewPersonService(srvc)
	err := handle.Update(context.Background(), &vladimir)
	assert.NoError(t, err)
	testVladimir, err := handle.ReadRow(context.Background(), vladimir.ID)
	assert.NoError(t, err)
	assert.Equal(t, testVladimir.ID, vladimir.ID)
	assert.Equal(t, testVladimir.Salary, vladimir.Salary)
	assert.Equal(t, testVladimir.Married, vladimir.Married)
	assert.Equal(t, testVladimir.Profession, vladimir.Profession)
	srvc.AssertExpectations(t)
}

// TestDelete is a mocktest for Delete method of interface Service
func TestDelete(t *testing.T) {
	srvc.On("Delete", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)

	handle := service.NewPersonService(srvc)
	err := handle.Delete(context.Background(), vladimir.ID)
	assert.NoError(t, err)

	srvc.AssertExpectations(t)
}
