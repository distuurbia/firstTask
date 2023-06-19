package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"os"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var rps *PersonPgx

var pgxVladimir = model.Person{
	Id:         uuid.New(),
	Salary:     2000,
	Married:    true,
	Profession: "policeman",
}

func PgxConnect() (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig("postgres://personuser:minovich12@localhost:5432/persondb")
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func TestMain(m *testing.M) {
	dbpool, err := PgxConnect()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
	}
	defer dbpool.Close()

	rps = NewRepository(dbpool)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://personUserMongoDB:minovich12@localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	mongoRps = NewMongoRep(client)
	exitVal := m.Run()
	os.Exit(exitVal)

}

func Test_PgxCreate(t *testing.T) {
	err := rps.Create(context.Background(), &pgxVladimir)
	require.NoError(t, err)
	testVladimir, err := rps.ReadRow(context.Background(), pgxVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, *testVladimir, pgxVladimir)
}

func Test_PgxCreateNil(t *testing.T) {
	err := rps.Create(context.Background(), nil)
	require.True(t, errors.Is(err, ErrNil))

}

func Test_PgxReadRow(t *testing.T) {
	testVladimir, err := rps.ReadRow(context.Background(), pgxVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, *testVladimir, pgxVladimir)
}

func Test_PgxReadRowNotFound(t *testing.T) {
	var id uuid.UUID
	_, err := rps.ReadRow(context.Background(), id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_PgxUpdate(t *testing.T) {
	pgxVladimir.Salary = 700
	err := rps.Update(context.Background(), &pgxVladimir)
	require.NoError(t, err)
	testVladimir, err := rps.ReadRow(context.Background(), pgxVladimir.Id)
	require.NoError(t, err)
	require.Equal(t, testVladimir.Id, pgxVladimir.Id)
	require.Equal(t, testVladimir.Salary, pgxVladimir.Salary)
	require.Equal(t, testVladimir.Married, pgxVladimir.Married)
	require.Equal(t, testVladimir.Profession, pgxVladimir.Profession)
}

func Test_PgxUpdateNotFound(t *testing.T) {
	var emptyEntity model.Person
	err := rps.Update(context.Background(), &emptyEntity)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_Delete(t *testing.T) {
	err := rps.Delete(context.Background(), pgxVladimir.Id)
	require.NoError(t, err)
	_, err = rps.ReadRow(context.Background(), pgxVladimir.Id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_DeleteNotFound(t *testing.T) {
	var id uuid.UUID
	err := rps.Delete(context.Background(), id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}
