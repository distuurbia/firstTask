package main

//import
import (
	"github.com/distuurbia/firstTask/internal/handler"
	"github.com/distuurbia/firstTask/internal/repository"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)
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

func main() {
	dbpool, err := Connect()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
	}
	defer dbpool.Close()
	persPgx := repository.NewRepository(dbpool)
	srv := service.NewService(persPgx)
	handl := handler.NewHandler(srv)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/persondb", handl.Create)
	e.GET("/persondb/:id", handl.ReadRow)
	e.PUT("persondb/:id", handl.Update)
	e.DELETE("persondb/:id", handl.Delete)
	e.Logger.Fatal(e.Start(":8080"))

}
