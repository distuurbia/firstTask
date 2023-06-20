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

func (mongoRps *PersonMongo) Create(ctx context.Context, pers *model.Person) error {
	if pers == nil {
		return ErrNil
	}
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	_, err := coll.InsertOne(ctx, pers)
	if err != nil{
		return fmt.Errorf("failed to create: %w", err)
	}
	return nil

}

func (mongoRps *PersonMongo) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": id}
	var pers model.Person
	err := coll.FindOne(ctx, filter).Decode(&pers)
	if err != nil {
		return &pers, fmt.Errorf("failed to read %w", err)
	}
	return &pers, nil
}

func (mongoRps *PersonMongo) GetAll(ctx context.Context) ([]model.Person, error) {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{}
	var allPers []model.Person
	cursor, err := coll.Find(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var pers model.Person
	for cursor.Next(ctx) {
		err = cursor.Decode(&pers)
		if err != nil {
			return allPers, fmt.Errorf("%w", err)
		}
		allPers = append(allPers, pers)
	}
	return allPers, nil
}

func (mongoRps *PersonMongo) Update(ctx context.Context, pers *model.Person) error {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": pers.Id}
	update := bson.M{"$set": pers}
	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (mongoRps *PersonMongo) Delete(ctx context.Context, id uuid.UUID) error {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": id}
	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

