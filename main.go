package main

import (
	"context"
	"fmt"
	"os"

	"github.com/distuurbia/firstTask/internal/model"
	"github.com/distuurbia/firstTask/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("postgres://personuser:minovich12@localhost:5432/perosndb"))
	if err != nil {
		fmt.Println(err, "dbpool")
		return
	}
	rps := repository.NewPerson(dbpool)
	var num int
	fmt.Println("Choose from 1.Create\n 2.Read\n 3.Update\n 4.Delete")
	fmt.Scan(&num)
	var vladimir model.Person
	vladimir.Salary = 200
	vladimir.Married = true
	vladimir.Profession = "doctor"
	var ctx context.Context
	switch num {
	case 1:
		rps.Create(ctx, &vladimir)
	case 2:
		fmt.Println("What is the Id of the row you wanna Read?")
		var rowId uuid.UUID
		fmt.Scan(&rowId)
		fmt.Println(rps.ReadRow(ctx, rowId))
	case 3:
		fmt.Println("What is the Id of the row you wanna Update?")
		var rowId uuid.UUID
		var married string
		fmt.Scan(&rowId)
		fmt.Println("Fill the info \n Salary:")
		fmt.Scan(&vladimir.Salary)
		fmt.Println("Married:")
		if married == "true" {
			vladimir.Married = true
		} else {
			vladimir.Married = false
		}
		fmt.Println("Profession:")
		fmt.Scan(&vladimir.Profession)
		rps.Update(ctx, &vladimir, rowId)
	case 4:
		fmt.Println("What is the Id of the row you wanna Delete?")
		var rowId uuid.UUID
		fmt.Scan(&rowId)
		rps.Delete(ctx, rowId)
	}
}
