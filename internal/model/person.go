// Package model contains structs that we will use as a records in our dbs
package model

import "github.com/google/uuid"

// Person struct - type of a record for dbs
type Person struct {
	ID         uuid.UUID `json:"ID" bson:"_id"`
	Salary     int
	Married    bool
	Profession string
}
