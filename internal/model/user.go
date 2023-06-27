// Package model contains structs that we will use as a records in our dbs
package model

import "github.com/google/uuid"

// User contains an info about the user and will be written in a users table
type User struct {
	ID           uuid.UUID `json:"ID" bson:"_id"`
	Username     string    `json:"username" bson:"username" validate:"required,min=4,max=15"`
	Password     []byte    `json:"password" bson:"password" validate:"required,min=4,max=15"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
}
