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
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var rps *PersonPgx

var pgxVladimir = model.Person{
	ID:         uuid.New(),
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
		dbpool.Close()
		os.Exit(1)
	}
	rps = NewPgxRep(dbpool)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://personUserMongoDB:minovich12@localhost:27017"))
	if err != nil {
		fmt.Println(err)
		cancel()
		os.Exit(1)
	}
	mongoRps = NewMongoRep(client)
	exitVal := m.Run()
	dbpool.Close()
	client.Disconnect(ctx)
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
