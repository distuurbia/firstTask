package model

import "github.com/google/uuid"

type Person struct {
	Id         uuid.UUID
	Salary     int
	Married    bool
	Profession string
}
