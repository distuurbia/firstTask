// Package config contains all system variables
package config

// Config contains system variables
type Config struct {
	SecretKey             string `env:"SECRET_KEY"`
	MongoConnectionString string `env:"MONGO_CONN_STRING"`
	PgxConnectionString   string `env:"PGX_CONN_STRING"`
}
