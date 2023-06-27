// Package model contains structs that we will use as a records in our dbs
package model

import "github.com/google/uuid"

// Person contains an info about the person and will be written in a personsdb table
type Person struct {
	ID         uuid.UUID `json:"ID" bson:"_id"`
	Salary     int       `json:"salary" bson:"salary" validate:"required,numeric,min=100,max=100000"`
	Married    bool      `json:"married" bson:"married"`
	Profession string    `json:"profession" bson:"profession" validate:"required,min=3,max=30"`
}
