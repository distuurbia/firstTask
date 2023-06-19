package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PersonMongo struct {
	client *mongo.Client
}

func NewMongoRep(client *mongo.Client) *PersonMongo {
	return &PersonMongo{client: client}
}

func (mongoRps *PersonMongo) CreateMongo(ctx context.Context, pers *model.Person) error {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	_, err := coll.InsertOne(ctx, pers)
	if err != nil{
		return fmt.Errorf("failed to create: %w", err)
	}
	return nil

}

func (mongoRps *PersonMongo) ReadRowMongo(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"id": id}
	var pers model.Person
	err := coll.FindOne(ctx, filter).Decode(&pers)
	if err != nil {
		return &pers, fmt.Errorf("failed to read %w", err)
	}
	return &pers, nil
}

