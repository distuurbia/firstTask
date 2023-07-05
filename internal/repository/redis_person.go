// Package repository is a package for work with db methods
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Redis contains an object of type *redis.Client
type Redis struct {
	client *redis.Client
}

// NewRepositoryRedis accepts an object of *redis.Client and returns an object of type *Redis
func NewRepositoryRedis(client *redis.Client) *Redis {
	return &Redis{client: client}
}

// Set sets cache of person in redis db
func (rds *Redis) Set(ctx context.Context, pers *model.Person) error {
	persJSON, err := json.Marshal(pers)
	if err != nil {
		return fmt.Errorf("Redis -> Set -> json.Marshal -> error: %w", err)
	}
	rds.client.HSet(ctx, "person", pers.ID.String(), persJSON)
	return nil
}

// Get gets cache of person from redis db
func (rds *Redis) Get(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	persJSON, err := rds.client.HGet(ctx, "person", id.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, fmt.Errorf("Redis -> Get -> client.HGet -> error: %w", err)
	}
	var pers model.Person
	err = json.Unmarshal([]byte(persJSON), &pers)
	if err != nil {
		return nil, fmt.Errorf("Redis -> Get -> json.Unmarshal -> error: %w", err)
	}
	return &pers, nil
}

// Delete deletes cache of person from redis db
func (rds *Redis) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := rds.client.HDel(ctx, "person", id.String()).Result()
	if err != nil {
		return fmt.Errorf("Redis -> Delete -> client.HDel -> error: %w", err)
	}
	return nil
}

// // AddToStream adds to person_stream person
// func (rds *Redis) AddToStream(ctx context.Context, msg string) error {
// 	streamData := redis.XAddArgs{
// 		Stream: "person_stream",
// 		Values: map[string]interface{}{
// 			"data": msg,
// 		},
// 	}
// 	_, err := rds.client.XAdd(ctx, &streamData).Result()
// 	if err != nil {
// 		return fmt.Errorf("Redis -> AddToStream -> XAdd -> error: %w", err)
// 	}
// 	return nil
// }

// // GetFromStream gets a person from the Redis Stream by ID
// func (rds *Redis) GetFromStream(ctx context.Context, id uuid.UUID) (string, error) {
// 	streamData := redis.XReadArgs{
// 		Streams: []string{"person_stream", "0"},
// 		Count:   1,
// 		Block:   0,
// 	}
// 	results, err := rds.client.XRead(ctx, &streamData).Result()
// 	if err != nil {
// 		return "", fmt.Errorf("Redis -> GetFromStream -> XRead -> error: %w", err)
// 	}
// 	if len(results) == 0 || len(results[0].Messages) == 0 {
// 		return "", redis.Nil
// 	}
// 	var pers model.Person
// 	var persJSON string
// 	for _, msg := range results[0].Messages {
// 		if msg.Values["MyID"] == id.String() {
// 			persJSON = msg.Values["data"].(string)
// 			err = json.Unmarshal([]byte(persJSON), &pers)
// 			if err != nil {
// 				return nil, fmt.Errorf("RedisStreamRepository -> GetFromStream -> json.Unmarshal -> error: %w", err)
// 			}
// 			return &pers, nil
// 		}
// 	}
// 	return nil, redis.Nil
// }

// // DeleteFromStream deletes a person from the Redis Stream by ID
// func (rds *Redis) DeleteFromStream(ctx context.Context, id uuid.UUID) error {
// 	streamData := redis.XReadArgs{
// 		Streams: []string{"person_stream", "0"},
// 		Count:   0,
// 		Block:   0,
// 	}
// 	results, err := rds.client.XRead(ctx, &streamData).Result()
// 	if err != nil {
// 		return fmt.Errorf("Redis -> GetFromStream -> XRead -> error: %w", err)
// 	}
// 	if len(results) == 0 || len(results[0].Messages) == 0 {
// 		return redis.Nil
// 	}
// 	var msgID string
// 	for _, msg := range results[0].Messages {
// 		if msg.Values["MyID"] == id.String() {
// 			msgID = msg.ID
// 		}
// 	}
// 	_, err = rds.client.XDel(ctx, "person_stream", msgID).Result()
// 	if err != nil {
// 		return fmt.Errorf("Redis -> DeleteFromStream -> XDel -> error: %w", err)
// 	}
// 	return nil
// }
