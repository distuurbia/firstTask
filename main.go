package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/distuurbia/firstTask/internal/handler"
	"github.com/distuurbia/firstTask/internal/repository"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
func ConnectPgx() (*pgxpool.Pool, error) {
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
func ConnectMongo()  (*mongo.Client, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://personUserMongoDB:minovich12@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	return client, nil
}
func main() {
	var handl *handler.PersonHandler
	fmt.Println("What db do u wanna use?\n 1.PostgreSQL\n 2.MongoDB")
	var dbChoose int
	fmt.Scan(&dbChoose)
	if dbChoose == 1{
		dbpool, err := ConnectPgx()
		if err != nil {
			log.Fatal("could not construct the pool: ", err)
		}
		defer dbpool.Close()
		persPgx := repository.NewPgxRep(dbpool)
		srv := service.NewService(persPgx)
		handl = handler.NewHandler(srv)
		
	} else if dbChoose == 2{
		client, err := ConnectMongo()
		if err != nil{
			log.Fatal("could not construct the client: ", err)
		}
		persMongo := repository.NewMongoRep(client)
		srv := service.NewService(persMongo)
		handl = handler.NewHandler(srv)
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				log.Fatal("%w", err)
			}
		}()
	} else {
		fmt.Println("The wrong number!")
		os.Exit(1)
	}

	
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/persondb", handl.Create)
	e.GET("/persondb/:id", handl.ReadRow)
	e.PUT("persondb/:id", handl.Update)
	e.DELETE("persondb/:id", handl.Delete)
	e.Logger.Fatal(e.Start(":8080"))

}
