// Package model contains structs that we will use as a records in our dbs
package model

import "github.com/google/uuid"

// Person contains an info about the person and will be written in a personsdb table
type Person struct {
	ID         uuid.UUID `json:"-" bson:"_id"`
	Salary     int       `json:"salary" bson:"salary" validate:"required,numeric,min=100,max=100000"`
	Married    bool      `json:"married" bson:"married"`
	Profession string    `json:"profession" bson:"profession" validate:"required,min=3,max=30"`
}

// User contains an info about the user and will be written in a users table
type User struct {
	ID           uuid.UUID `json:"-" bson:"_id"`
	Username     string    `json:"username" bson:"username" validate:"required,min=4,max=15"`
	Password     []byte    `json:"password" bson:"password" validate:"required,min=4,max=15"`
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
}

// UserRequest contains request for user binding
type UserRequest struct {
	Username string `json:"username" bson:"username" validate:"required,min=4,max=15"`
	Password string `json:"password" bson:"password" validate:"required,min=4,max=15"`
}

// RefreshRequest contains request for user refresh method
type RefreshRequest struct {
	AccessToken  string `json:"accessToken" bson:"accessToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}
