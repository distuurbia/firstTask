package main

// import (
//     "context"
// 	"fmt"
// 	//"os"

// 	"github.com/distuurbia/firstTask/internal/model"
// 	"github.com/distuurbia/firstTask/internal/repository"
// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// func postgresql(connStr string, maxConns int32) (*pgxpool.Pool, error){
// 	ctx := context.Background()
// 	cfg, err := pgxpool.ParseConfig(connStr)
// 	if err != nil{
// 		return nil, fmt.Errorf("failed to parse postgresql connection config: %w", err)
// 	}

// 	cfg.MaxConns = maxConns

// 	pool, err := pgxpool.NewWithConfig(ctx, cfg)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to establish connection to postgresql: %w", err)
// 	}
// 	err = pool.Ping(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get response from postgresql: %w", err)
// 	}
// 	return pool, nil
// }

func main() {
	// cfg, err := pgxpool.ParseConfig("postgres://personuser:minovich12@localhost:5432/persondb")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// rps := repository.NewPerson(dbpool)
	// var num int
	// fmt.Println("Choose from 1.Create \n 2.Read \n 3.Update \n 4.Delete")
	// fmt.Scan(&num)
	// var vladimir model.Person
	// vladimir.Salary = 200
	// vladimir.Married = true
	// vladimir.Profession = "doctor"
	// ctx := context.Background()
	// switch num {
	// case 1:
	// 	rps.Create(ctx, &vladimir)
	// // case 2:
	// // 	fmt.Println("What is the Id of the row you wanna Read?")
	// // 	var sal int
	// // 	fmt.Scan(&sal)
	// // 	if err != nil {
	// // 		fmt.Println(err, "parse rowid")
	// // 		return
	// // 	}
	// // 	fmt.Println(rps.ReadRow(ctx, sal))
	// case 3:
	// 	fmt.Println("What is the Id of the row you wanna Update?")
	// 	var rowId uuid.UUID
	// 	var married string
	// 	fmt.Scan(&rowId)
	// 	fmt.Println("Fill the info \n Salary:")
	// 	fmt.Scan(&vladimir.Salary)
	// 	fmt.Println("Married:")
	// 	if married == "true" {
	// 		vladimir.Married = true
	// 	} else {
	// 		vladimir.Married = false
	// 	}
	// 	fmt.Println("Profession:")
	// 	fmt.Scan(&vladimir.Profession)
	// 	rps.Update(ctx, &vladimir, rowId)
	// case 4:
	// 	fmt.Println("What is the Id of the row you wanna Delete?")
	// 	var rowId uuid.UUID
	// 	fmt.Scan(&rowId)
	// 	rps.Delete(ctx, rowId)
	// }
}
