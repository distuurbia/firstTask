// Package repository is a package for work with db methods
package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// PersonMongo contains object of type *mongo.Client
type PersonMongo struct {
	client *mongo.Client
}

// NewMongoRep accepts object of type *mongo.Client and returns an object of type *PersonMongo
func NewMongoRep(client *mongo.Client) *PersonMongo {
	return &PersonMongo{client: client}
}

// Create creates document in mongoDB collection
func (mongoRps *PersonMongo) Create(ctx context.Context, pers *model.Person) error {
	if pers == nil {
		return ErrNil
	}
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	_, err := coll.InsertOne(ctx, pers)
	if err != nil {
		return fmt.Errorf("failed to create: %w", err)
	}
	return nil
}

// ReadRow reads document from mongoDB collection
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

// GetAll reads all documents from mongoDB collection
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

// Update update the document of mongoDB collection
func (mongoRps *PersonMongo) Update(ctx context.Context, pers *model.Person) error {
	coll := mongoRps.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": pers.ID}
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

// Delete deletes the document of mongoDB collection
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
