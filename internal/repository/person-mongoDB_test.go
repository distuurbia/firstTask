package repository

import (
	"context"
	// "log"
	"testing"
	// "time"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)
var mongoRps *PersonMongo

var mongoVladimir = model.Person{
	Id: uuid.New(),
	Salary: 2000,
	Married: true,
	Profession: "policeman",
}


func Test_MongoCreate(t *testing.T){
	// ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
	// defer cancel()
	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://personUserMongoDB:minovich12@localhost:27017"))


	// if err != nil {
	// 	log.Fatal("Could not construct the pool: ", err)
	// }
	// mongoRps = NewMongoRep(client)
	err := mongoRps.CreateMongo(context.Background(), &mongoVladimir)
	require.NoError(t, err)
	testMongoVladimir, err := mongoRps.ReadRowMongo(context.Background(), mongoVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.Id, testMongoVladimir.Id)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}