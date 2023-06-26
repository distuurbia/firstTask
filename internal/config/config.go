// Package config contains all system variables
package config

const (
	// SecretKey is my secret key for signing jwt tokens
	SecretKey = "roapjumping"
	// MongoConnectionString is mongo connection string
	MongoConnectionString = "mongodb://personUserMongoDB:minovich12@localhost:27017"
	// PgxConnectionString is pgx connection string
	PgxConnectionString = "postgres://personuser:minovich12@localhost:5432/persondb"
)
