package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var rps *PersonPgx

var vladimir = model.Person{
	Id:         uuid.New(),
	Salary:     2000,
	Married:    true,
	Profession: "policeman",
}

func Connect() (*pgxpool.Pool, error) {
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
	dbpool, err := Connect()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
	}
	defer dbpool.Close()

	rps = NewRepository(dbpool)
	exitVal := m.Run()
	os.Exit(exitVal)

}

func Test_Create(t *testing.T) {
	err := rps.Create(context.Background(), &vladimir)
	require.NoError(t, err)
	var testVladimir *model.Person
	testVladimir, err = rps.ReadRow(context.Background(), vladimir.Id)
	require.NoError(t, err)
	require.Equal(t, *testVladimir, vladimir)
}

func Test_CreateNil(t *testing.T) {
	err := rps.Create(context.Background(), nil)
	require.True(t, errors.Is(err, repository.ErrNil))

}

func Test_ReadRow(t *testing.T) {
	testVladimir, err := rps.ReadRow(context.Background(), vladimir.Id)
	require.NoError(t, err)
	require.Equal(t, *testVladimir, vladimir)
}

func Test_ReadRowNotFound(t *testing.T) {
	var id uuid.UUID
	_, err := rps.ReadRow(context.Background(), id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_Update(t *testing.T) {
	vladimir.Salary = 700
	err := rps.Update(context.Background(), &vladimir)
	require.NoError(t, err)
	testVladimir, err := rps.ReadRow(context.Background(), vladimir.Id)

}

func Test_UpdateNotFound(t *testing.T) {
	var emptyEntity model.Person
	err := rps.Update(context.Background(), &emptyEntity)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_Delete(t *testing.T) {
	err := rps.Delete(context.Background(), vladimir.Id)
	require.NoError(t, err)
	_, err = rps.ReadRow(context.Background(), vladimir.Id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}

func Test_DeleteNotFound(t *testing.T) {
	var id uuid.UUID
	err := rps.Delete(context.Background(), id)
	require.True(t, errors.Is(err, pgx.ErrNoRows))
}
