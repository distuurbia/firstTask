package repository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var mongoRps *PersonMongo

var vladimir = model.Person{
	Id: uuid.New(),
	Salary: 2000,
	Married: true,
	Profession: "policeman",
}

func TestMain(m *testing.M){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Could not construct the pool: ", err)
	}
	mongoRps = NewMongoRep(client)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_Create(t *testing.T){
	mongoRps.Create(context.Background(), &vladimir)
}