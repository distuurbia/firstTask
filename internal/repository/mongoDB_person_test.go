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

var rpsMongo *Mongo

var mongoVladimir = model.Person{
	ID:         uuid.New(),
	Salary:     2000,
	Married:    true,
	Profession: "policeman",
}

func Test_MongoCreate(t *testing.T) {
	err := rpsMongo.Create(context.Background(), &mongoVladimir)
	require.NoError(t, err)
	testMongoVladimir, err := rpsMongo.ReadRow(context.Background(), mongoVladimir.ID)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.ID, testMongoVladimir.ID)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}

func Test_MongoCreateNil(t *testing.T) {
	err := rpsMongo.Create(context.Background(), nil)
	require.True(t, errors.Is(err, ErrNil))
}

func Test_MongoCreateDuplicate(t *testing.T) {
	err := rpsMongo.Create(context.Background(), &mongoVladimir)
	require.Error(t, err)
}

func Test_MongoCreateContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := rpsMongo.Create(ctx, &mongoVladimir)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_MongoReadRow(t *testing.T) {
	testMongoVladimir, err := rpsMongo.ReadRow(context.Background(), mongoVladimir.ID)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.ID, testMongoVladimir.ID)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}

func Test_MongoReadRowNotFound(t *testing.T) {
	var id uuid.UUID
	_, err := rpsMongo.ReadRow(context.Background(), id)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoReadRowContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	_, err := rpsMongo.ReadRow(ctx, mongoVladimir.ID)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_MongoGetAll(t *testing.T) {
	allPers, err := rpsMongo.GetAll(context.Background())
	require.NoError(t, err)
	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{}
	numberPersons, err := coll.CountDocuments(context.Background(), filter)
	require.NoError(t, err)
	require.Equal(t, len(allPers), int(numberPersons))
}

func Test_MongoUpdate(t *testing.T) {
	mongoVladimir.Salary = 100
	mongoVladimir.Married = false
	mongoVladimir.Profession = "Security"
	err := rpsMongo.Update(context.Background(), &mongoVladimir)
	require.NoError(t, err)
	testMongoVladimir, err := rpsMongo.ReadRow(context.Background(), mongoVladimir.ID)
	require.NoError(t, err)
	require.Equal(t, mongoVladimir.ID, testMongoVladimir.ID)
	require.Equal(t, mongoVladimir.Salary, testMongoVladimir.Salary)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
	require.Equal(t, mongoVladimir.Married, testMongoVladimir.Married)
}

func Test_MongoUpdateNotFound(t *testing.T) {
	var emptyEntity model.Person
	err := rpsMongo.Update(context.Background(), &emptyEntity)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoUpdateContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := rpsMongo.Update(ctx, &mongoVladimir)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_MongoDelete(t *testing.T) {
	err := rpsMongo.Delete(context.Background(), mongoVladimir.ID)
	require.NoError(t, err)
	_, err = rpsMongo.ReadRow(context.Background(), mongoVladimir.ID)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoDeleteNotFound(t *testing.T) {
	var id uuid.UUID
	err := rpsMongo.Delete(context.Background(), id)
	require.True(t, errors.Is(err, mongo.ErrNoDocuments))
}

func Test_MongoDeleteontextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := rpsMongo.Delete(ctx, mongoVladimir.ID)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}
