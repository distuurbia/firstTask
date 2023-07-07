package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const(
	pgUsername = "personuser"
	pgPassword = "minovich12"
	pgDB = "persondb"
)

var rps *Pgx

var pgxVladimir = model.Person{
	ID:         uuid.New(),
	Salary:     2000,
	Married:    true,
	Profession: "policeman",
}

func SetupTestPgx() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		fmt.Sprintf("POSTGRES_USER=%s", pgUsername),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", pgPassword),
		fmt.Sprintf("POSTGRES_DB=%s", pgDB)})
	if err != nil {
		logrus.Fatalf("can't start postgres container: %s", err)
	}
	cmd := exec.Command(
		"flyway",
		fmt.Sprintf("-user=%s", pgUsername),
		fmt.Sprintf("-password=%s", pgPassword),
		"-locations=filesystem:../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%s:%s/persondb", "localhost", resource.GetPort("5432/tcp")), "-connectRetries=10",
		"migrate",
	)
	err = cmd.Run()
	if err != nil {
		logrus.Fatalf("can't run migration: %s", err)
	}
	port := resource.GetPort("5432/tcp")
	dbURL := fmt.Sprintf("postgresql://personuser:minovich12@localhost:%s/persondb", port)
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logrus.Fatalf("can't parse config: %s", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logrus.Fatalf("can't connect to postgtres: %s", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func SetupTestMongoDB() (*mongo.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("mongo", "latest", []string{
		"MONGO_INITDB_ROOT_USERNAME=personUserMongoDB",
		"MONGO_INITDB_ROOT_PASSWORD=minovich12",
		"MONGO_INITDB_DATABASE=personMongoDB"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	port := resource.GetPort("27017/tcp")
	mongoURL := fmt.Sprintf("mongodb://personUserMongoDB:minovich12@localhost:%s", port)
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect mongoDB: %w", err)
	}
	cleanup := func() {
		client.Disconnect(context.Background())
		pool.Purge(resource)
	}
	return client, cleanup, nil
}

func TestMain(m *testing.M) {
	dbpool, cleanupPgx, err := SetupTestPgx()
	if err != nil {
		cleanupPgx()
		os.Exit(1)
	}
	rps = NewRepositoryPgx(dbpool)
	client, cleanupMongo, err := SetupTestMongoDB()
	if err != nil {
		fmt.Println(err)
		cleanupMongo()
		os.Exit(1)
	}
	rpsMongo = NewRepositoryMongo(client)
	exitVal := m.Run()
	cleanupPgx()
	cleanupMongo()
	os.Exit(exitVal)
}
