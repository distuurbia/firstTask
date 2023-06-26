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

// Mongo contains object of type *mongo.Client
type Mongo struct {
	client *mongo.Client
}

// NewRepositoryMongo accepts object of type *mongo.Client and returns an object of type *PersonMongo
func NewRepositoryMongo(client *mongo.Client) *Mongo {
	return &Mongo{client: client}
}

// Create creates document in mongoDB collection
func (rpsMongo *Mongo) Create(ctx context.Context, pers *model.Person) error {
	if pers == nil {
		return ErrNil
	}
	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")
	_, err := coll.InsertOne(ctx, pers)
	if err != nil {
		return fmt.Errorf("PersonMongo -> Create -> error: %w", err)
	}
	return nil
}

// ReadRow reads document from mongoDB collection
func (rpsMongo *Mongo) ReadRow(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": id}
	var pers model.Person
	err := coll.FindOne(ctx, filter).Decode(&pers)
	if err != nil {
		return &pers, fmt.Errorf("PersonMongo -> ReadRow -> error: %w", err)
	}
	return &pers, nil
}

// GetAll reads all documents from mongoDB collection
func (rpsMongo *Mongo) GetAll(ctx context.Context) ([]model.Person, error) {
	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{}
	var allPers []model.Person
	cursor, err := coll.Find(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("PersonMongo -> GetAll -> Find -> error: %w", err)
	}
	var pers model.Person
	for cursor.Next(ctx) {
		err = cursor.Decode(&pers)
		if err != nil {
			return allPers, fmt.Errorf("PersonMongo -> GetAll -> Decode -> error: %w", err)
		}
		allPers = append(allPers, pers)
	}
	return allPers, nil
}

// Update update the document of mongoDB collection
func (rpsMongo *Mongo) Update(ctx context.Context, pers *model.Person) error {
	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": pers.ID}
	update := bson.M{"$set": pers}
	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("PersonMongo -> Update -> error: %w", err)
	}
	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Delete deletes the document of mongoDB collection
func (rpsMongo *Mongo) Delete(ctx context.Context, id uuid.UUID) error {
	coll := rpsMongo.client.Database("personMongoDB").Collection("persons")
	filter := bson.M{"_id": id}
	res, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("PersonMongo -> Delete -> error: %w", err)
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
