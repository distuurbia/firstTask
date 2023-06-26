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

// SignUp create new user in users collection
func (rpsMongo *Mongo) SignUp(ctx context.Context, user *model.User) error {
	if user == nil {
		return ErrNil
	}
	coll := rpsMongo.client.Database("personMongoDB").Collection("users")
	filter := bson.M{"username": user.Username}
	numberPersons, err := coll.CountDocuments(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("Mongo -> SignUp -> CountDocuments -> error: %w", err)
	}
	if numberPersons != 0 {
		return ErrExist
	}
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("Mongo -> SignUp -> InsertOne -> error: %w", err)
	}
	return nil
}

// GetPasswordAndIDByUsername returnes id and hash of the password from users table
func (rpsMongo *Mongo) GetPasswordAndIDByUsername(ctx context.Context, username string) (uuid.UUID, []byte, error) {
	var user model.User
	coll := rpsMongo.client.Database("personMongoDB").Collection("users")
	filter := bson.M{"username": username}
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return uuid.UUID{}, nil, fmt.Errorf("Mongo -> GetPasswordAndIDByUserName -> FindOne -> error: %w", err)
	}
	return user.ID, user.Password, nil
}

// GetRefreshTokenByID returnes refreshToken from users table by id
func (rpsMongo *Mongo) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error) {
	coll := rpsMongo.client.Database("personMongoDB").Collection("users")
	filter := bson.M{"_id": id}
	var user *model.User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("Mongo -> GetRefreshTokenByName -> QueryRow -> error: %w", err)
	}
	return user.RefreshToken, nil
}

// AddRefreshToken adds refreshToken to users table by id
func (rpsMongo *Mongo) AddRefreshToken(ctx context.Context, user *model.User) error {
	coll := rpsMongo.client.Database("personMongoDB").Collection("users")
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{"refreshToken": user.RefreshToken}}
	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Mongo -> AddRefreshToken -> UpdateOne -> error: %w", err)
	}
	if res.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
