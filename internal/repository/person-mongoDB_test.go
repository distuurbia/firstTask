package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var mongoRps *PersonMongo

var mongoVladimir = model.Person{
	Id:         uuid.New(),
	Salary:     2000,
	Married:    true,
	Profession: "policeman",
}

func Test_MongoCreate(t *testing.T) {
	err := mongoRps.Create(context.Background(), &mongoVladimir)
	require.NoError(t, err)
	testMongoVladimir, err := mongoRps.ReadRow(context.Background(), mongoVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.Id, testMongoVladimir.Id)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}

func Test_MongoCreateNil(t *testing.T) {
	err := mongoRps.Create(context.Background(), nil)
	require.True(t, errors.Is(err, ErrNil))

}

func Test_MongoCreateDuplicate(t *testing.T) {
	err := mongoRps.Create(context.Background(), &mongoVladimir)
	require.Error(t, err)
}

func Test_MongoCreateContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := mongoRps.Create(ctx, &mongoVladimir)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_MongoReadRow(t *testing.T) {
	testMongoVladimir, err := mongoRps.ReadRow(context.Background(), mongoVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.Id, testMongoVladimir.Id)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}

func Test_MongoReadRowNotFound(t *testing.T) {
	var id uuid.UUID
	_, err := mongoRps.ReadRow(context.Background(), id)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoReadRowContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	_, err := mongoRps.ReadRow(ctx, mongoVladimir.Id)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_MongoGetAll(t *testing.T) {
	allPers, err := mongoRps.GetAll(context.Background())
	require.NoError(t, err)
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")

	filter := bson.M{}

	numberPersons, err := coll.CountDocuments(context.Background(), filter)

	require.NoError(t, err)
	require.Equal(t, len(allPers), int(numberPersons))

}

func Test_MongoUpdate(t *testing.T) {
	mongoVladimir.Salary = 100
	mongoVladimir.Married = false
	mongoVladimir.Profession = "Security"
	err := mongoRps.Update(context.Background(), &mongoVladimir)
	require.NoError(t, err)
	testMongoVladimir, err := mongoRps.ReadRow(context.Background(), mongoVladimir.Id)
	require.Equal(t, mongoVladimir.Id, testMongoVladimir.Id)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}

func Test_MongoUpdateNotFound(t *testing.T) {
	var emptyEntity model.Person
	err := mongoRps.Update(context.Background(), &emptyEntity)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoUpdateContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := mongoRps.Update(ctx, &mongoVladimir)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_MongoDelete(t *testing.T) {
	err := mongoRps.Delete(context.Background(), mongoVladimir.Id)
	require.NoError(t, err)
	_, err = mongoRps.ReadRow(context.Background(), mongoVladimir.Id)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoDeleteNotFound(t *testing.T) {
	var id uuid.UUID
	err := mongoRps.Delete(context.Background(), id)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoDeleteontextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := mongoRps.Delete(ctx, mongoVladimir.Id)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}
