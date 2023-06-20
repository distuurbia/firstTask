package handler

import (
	"context"
	"os"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
 	"github.com/distuurbia/firstTask/internal/handler/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var vladimir = model.Person{
	Id: uuid.New(),
	Salary: 300,
	Married: false,
	Profession: "teacher",
}
var srvc *mocks.Service

func TestMain(m *testing.M){
	srvc = new(mocks.Service)
	exitVal := m.Run()
	os.Exit(exitVal)
}
func TestCreate(t *testing.T){
	srvc.On("Create", mock.Anything, mock.AnythingOfType("*model.Person")).Return(nil).Once()

	handl := service.NewService(srvc)
	err := handl.Create(context.Background(), &vladimir)
	assert.Nil(t, err)

	srvc.AssertExpectations(t)
}

func TestReadRow(t *testing.T){
	srvc.On("ReadRow", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&vladimir, nil)

	handle := service.NewService(srvc)
	testVladimir, err := handle.ReadRow(context.Background(), vladimir.Id)
	assert.Nil(t, err)
	assert.Equal(t, testVladimir.Id, vladimir.Id)
	assert.Equal(t, testVladimir.Salary, vladimir.Salary)
	assert.Equal(t, testVladimir.Married, vladimir.Married)
	assert.Equal(t, testVladimir.Profession, vladimir.Profession)

	srvc.AssertExpectations(t)
}

func TestUpdateRow(t *testing.T){
	srvc.On("Update", mock.Anything, mock.AnythingOfType("*model.Person")).Return(nil).Once()
	vladimir.Salary = 700
	vladimir.Married = false
	vladimir.Profession = "Lawer"
	handle := service.NewService(srvc)
	err := handle.Update(context.Background(), &vladimir)
	assert.Nil(t, err)
	testVladimir, err := handle.ReadRow(context.Background(), vladimir.Id)
	assert.Equal(t, testVladimir.Id, vladimir.Id)
	assert.Equal(t, testVladimir.Salary, vladimir.Salary)
	assert.Equal(t, testVladimir.Married, vladimir.Married)
	assert.Equal(t, testVladimir.Profession, vladimir.Profession)

	srvc.AssertExpectations(t)
}

func TestDelete(t *testing.T){
	srvc.On("Delete", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)

	handle := service.NewService(srvc)
	err := handle.Delete(context.Background(), vladimir.Id)
	assert.Nil(t, err)

	srvc.AssertExpectations(t)
}