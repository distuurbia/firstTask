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
	userJSON, err := json.Marshal(pers)
	if err != nil {
		return fmt.Errorf("Redis -> Set -> json.Marshal -> error: %w", err)
	}
	rds.client.HSet(ctx, "person", pers.ID.String(), userJSON)
	return nil
}

// Get gets cache of person from redis db
func (rds *Redis) Get(ctx context.Context, id uuid.UUID) (*model.Person, error) {
	userJSON, err := rds.client.HGet(ctx, "person", id.String()).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, err
		}
		return nil, fmt.Errorf("Redis -> Get -> client.HGet -> error: %w", err)
	}
	var pers model.Person
	err = json.Unmarshal([]byte(userJSON), &pers)
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
