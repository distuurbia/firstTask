package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
)

type entity interface {
	create(dbPointer **sql.DB)
	read()
	update()
	delete()
}

type person struct {
	//id         int
	salary     int
	married    bool
	profession string
}

func (p person) create(dbPointer **sql.DB) {
	stmt, err := (*dbPointer).Prepare("INSERT INTO persondb(salary, married, profession) VALUES($1, $2, $3)")
	if err != nil {
		panic(err)
	}
	stmt.Exec(p.salary, p.married, p.profession)
	defer stmt.Close()
}

func (p person) read() {

}

func (p person) update() {

}

func (p person) delete() {

}

func main() {
	connStr := "user=personuser password=minovich12 host=localhost port=5432 database=persondb sslmode=disable" // "pgx", "user=postgres password=secret host=localhost port=5432 database=pgx_test sslmode=disable"

	fmt.Println(connStr)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var vladimir entity = person{salary: 200, married: true, profession: "doctor"}
	vladimir.create(&db)
}
