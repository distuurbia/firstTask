package repository

import (
	"context"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)
var mongoRps *PersonMongo

var mongoVladimir = model.Person{
	Id: uuid.New(),
	Salary: 2000,
	Married: true,
	Profession: "policeman",
}


func Test_MongoCreate(t *testing.T){
	err := mongoRps.CreateMongo(context.Background(), &mongoVladimir)
	require.NoError(t, err)
	testMongoVladimir, err := mongoRps.ReadRowMongo(context.Background(), mongoVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.Id, testMongoVladimir.Id)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}