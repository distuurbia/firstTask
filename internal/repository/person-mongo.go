package repository

import (
	"context"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
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
	coll := mongoRps.client.Database("mongoPerson").Collection("persons")
	_, err := coll.InsertOne(ctx, pers)
	return fmt.Errorf("failed to create: %w", err)

}

func (mongoRps *PersonMongo) ReadRow(ctx context.Context, pers *model.Person) (*model.Person, error) {
	coll := mongoRps.client.Database("mongoPerson").Collection("persons")
	filter := bson.M{"id": pers.Id}

	err := coll.FindOne(ctx, filter).Decode(&pers)
	if err != nil {

	}
	return pers, nil
}

// func ConnectMongo(){
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

// 	defer func() {
// 		if err = client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()

// }
