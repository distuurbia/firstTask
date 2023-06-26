package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"POSTGRES_USER=personuser",
		"POSTGRESQL_PASSWORD=minovich12",
		"POSTGRES_DB=persondb"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	dbURL := fmt.Sprintf("postgres://personuser:minovich12@localhost:%s/persondb", resource.GetPort("5432"))
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
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
		fmt.Println("Could not construct the pool: ", err)
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

func Test_PgxCreate(t *testing.T) {
	err := rps.Create(context.Background(), &pgxVladimir)
	require.NoError(t, err)
	testVladimir, err := rps.ReadRow(context.Background(), pgxVladimir.ID)
	require.NoError(t, err)
	require.Equal(t, testVladimir.ID, pgxVladimir.ID)
	require.Equal(t, testVladimir.Salary, pgxVladimir.Salary)
	require.Equal(t, testVladimir.Married, pgxVladimir.Married)
	require.Equal(t, testVladimir.Profession, pgxVladimir.Profession)
}

func Test_PgxCreateNil(t *testing.T) {
	err := rps.Create(context.Background(), nil)
	require.True(t, errors.Is(err, ErrNil))
}

func Test_PgxCreateDuplicate(t *testing.T) {
	err := rps.Create(context.Background(), &pgxVladimir)
	require.Error(t, err)
}

func Test_PgxCreateContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := rps.Create(ctx, &pgxVladimir)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}
func Test_PgxReadRow(t *testing.T) {
	testVladimir, err := rps.ReadRow(context.Background(), pgxVladimir.ID)
	require.NoError(t, err)
	require.Equal(t, testVladimir.ID, pgxVladimir.ID)
	require.Equal(t, testVladimir.Salary, pgxVladimir.Salary)
	require.Equal(t, testVladimir.Married, pgxVladimir.Married)
	require.Equal(t, testVladimir.Profession, pgxVladimir.Profession)
}

func Test_PgxReadRowNotFound(t *testing.T) {
	var id uuid.UUID
	_, err := rps.ReadRow(context.Background(), id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_PgxReadRowContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	_, err := rps.ReadRow(ctx, pgxVladimir.ID)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_PgxGetAll(t *testing.T) {
	allPers, err := rps.GetAll(context.Background())
	require.NoError(t, err)
	var numberPersons int
	err = rps.db.QueryRow(context.Background(), "SELECT COUNT(*) FROM persondb").Scan(&numberPersons)
	require.NoError(t, err)
	require.Equal(t, len(allPers), numberPersons)
}

func Test_PgxUpdate(t *testing.T) {
	pgxVladimir.Salary = 700
	pgxVladimir.Married = false
	pgxVladimir.Profession = "Lawer"
	err := rps.Update(context.Background(), &pgxVladimir)
	require.NoError(t, err)
	testVladimir, err := rps.ReadRow(context.Background(), pgxVladimir.ID)
	require.NoError(t, err)
	require.Equal(t, testVladimir.ID, pgxVladimir.ID)
	require.Equal(t, testVladimir.Salary, pgxVladimir.Salary)
	require.Equal(t, testVladimir.Married, pgxVladimir.Married)
	require.Equal(t, testVladimir.Profession, pgxVladimir.Profession)
}

func Test_PgxUpdateNotFound(t *testing.T) {
	var emptyEntity model.Person
	err := rps.Update(context.Background(), &emptyEntity)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_PgxUpdateContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := rps.Update(ctx, &pgxVladimir)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}

func Test_PgxDelete(t *testing.T) {
	err := rps.Delete(context.Background(), pgxVladimir.ID)
	require.NoError(t, err)
	_, err = rps.ReadRow(context.Background(), pgxVladimir.ID)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_PgxDeleteNotFound(t *testing.T) {
	var id uuid.UUID
	err := rps.Delete(context.Background(), id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_PgxDeleteContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	time.Sleep(1 * time.Second)
	defer cancel()
	err := rps.Delete(ctx, pgxVladimir.ID)
	require.True(t, errors.Is(err, context.DeadlineExceeded))
}
