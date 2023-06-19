package model

import "github.com/google/uuid"

type Person struct {
	Id         uuid.UUID `json:"Id" bson:"_id"`
	Salary     int
	Married    bool
	Profession string
}
