package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/distuurbia/firstTask/internal/model"

	//"github.com/google/uuid"

	// "github.com/distuurbia/firstTask/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var rps *Person

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
		fmt.Errorf("Could not construct the pool: %w", err)
	}
	defer dbpool.Close()

	rps = NewPerson(dbpool)
	exitVal := m.Run()
	os.Exit(exitVal)

}

func Test_Create(t *testing.T) {
	var vladimir model.Person

	vladimir.Salary = 3300
	vladimir.Married = true
	vladimir.Profession = "doctor"
	id, err := rps.Create(context.Background(), &vladimir)
	require.NoError(t, err)

	vladimir.Id = id

	var testVladimir *model.Person
	testVladimir, err = rps.ReadRow(context.Background(), id) //.Scan(&testVladimir.Id, &testVladimir.Salary, &testVladimir.Married, &testVladimir.Profession)
	require.NoError(t, err)
	require.Equal(t, *testVladimir, vladimir)
	// require.Equal(t, testVladimir.Salary, vladimir.Salary)
	// require.Equal(t, testVladimir.Married, vladimir.Married)
	// require.Equal(t, testVladimir.Profession, vladimir.Profession)
}

// func Test_ReadRow(t *testing.T) {
// 	dbpool, err := Connect()
// 	require.NoError(t, err)
// 	defer dbpool.Close()

// 	rps := NewPerson(dbpool)
// 	ctx := context.Background()
// 	vladimir := model.Person{Salary: 3300, Married: true, Profession: "doctor"}
// 	id := uuid.MustParse("686bd87a-94ee-4c82-94e5-d82618a9ced6")
// 	var testVladimir model.Person
// 	testVladimir, err := rps.ReadRow(ctx, id)
// 	require.NoError(t, err)

// 	vladimir.Id = id

// 	err = dbpool.QueryRow(ctx, "SELECT id, salary, married, profession FROM persondb WHERE id = $1", id).Scan(&testVladimir.Id, &testVladimir.Salary, &testVladimir.Married, &testVladimir.Profession)
// 	require.NoError(t, err)
// 	require.Equal(t, testVladimir, vladimir)
// }
