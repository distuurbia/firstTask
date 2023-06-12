package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx"
)

func main() {
	connStr := "user=personuser password=minovich12 dbname=persondb sslmode=disable"
	fmt.Println(connStr)
	db, err := sql.Open("postgress", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

}
