package handler

import (
	"context"
	"os"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/distuurbia/firstTask/internal/service/mocks"
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
var rps *mocks.Repository

func TestMain(m *testing.M){
	rps = new(mocks.Repository)
	exitVal := m.Run()
	os.Exit(exitVal)
}
func TestCreate(t *testing.T){
	rps.On("Create", mock.Anything, mock.AnythingOfType("*model.Person")).Return(nil).Once()

	srv := service.NewService(rps)
	err := srv.Create(context.Background(), &vladimir)
	assert.Nil(t, err)

	rps.AssertExpectations(t)
}

func TestReadRow(t *testing.T){
	rps.On("ReadRow", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(&vladimir, nil)

	srv := service.NewService(rps)
	testVladimir, err := srv.ReadRow(context.Background(), vladimir.Id)
	assert.Nil(t, err)
	assert.Equal(t, testVladimir.Id, vladimir.Id)
	assert.Equal(t, testVladimir.Salary, vladimir.Salary)
	assert.Equal(t, testVladimir.Married, vladimir.Married)
	assert.Equal(t, testVladimir.Profession, vladimir.Profession)

	rps.AssertExpectations(t)
}

func TestUpdateRow(t *testing.T){
	rps.On("Update", mock.Anything, mock.AnythingOfType("*model.Person")).Return(nil).Once()
	vladimir.Salary = 700
	vladimir.Married = false
	vladimir.Profession = "Lawer"
	srv := service.NewService(rps)
	err := srv.Update(context.Background(), &vladimir)
	assert.Nil(t, err)
	testVladimir, err := srv.ReadRow(context.Background(), vladimir.Id)
	assert.Equal(t, testVladimir.Id, vladimir.Id)
	assert.Equal(t, testVladimir.Salary, vladimir.Salary)
	assert.Equal(t, testVladimir.Married, vladimir.Married)
	assert.Equal(t, testVladimir.Profession, vladimir.Profession)

	rps.AssertExpectations(t)
}

func TestDelete(t *testing.T){
	rps.On("Delete", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)

	srv := service.NewService(rps)
	err := srv.Delete(context.Background(), vladimir.Id)
	assert.Nil(t, err)

	rps.AssertExpectations(t)
}