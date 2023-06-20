// Package main contains main function
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

// ConnectPgx connects to the pgxpool
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

// ConnectMongo connects to the mongoDB
func ConnectMongo() (*mongo.Client, error) {
	const ctxTimeout = 10
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://personUserMongoDB:minovich12@localhost:27017"))
	if err != nil {
		return client, fmt.Errorf("%w", err)
	}
	return client, nil
}

// main is an executable function
func main() {
	var handl *handler.PersonHandler
	fmt.Println("What db do u wanna use?\n 1.PostgreSQL\n 2.MongoDB")
	var dbChoose int
	dbChoose, err := fmt.Scan()
	if err != nil {
		fmt.Println("failed to scan")
	}
	const PostgreSQL = 1
	const MongoDB = 2
	switch dbChoose {
	case PostgreSQL:
		dbpool, err := ConnectPgx()
		if err != nil {
			log.Fatal("could not construct the pool: ", err)
		}
		defer dbpool.Close()
		persPgx := repository.NewPgxRep(dbpool)
		srv := service.NewService(persPgx)
		handl = handler.NewHandler(srv)
	case MongoDB:
		client, err := ConnectMongo()
		if err != nil {
			fmt.Println("could not construct the client: ", err)
		}
		persMongo := repository.NewMongoRep(client)
		srv := service.NewService(persMongo)
		handl = handler.NewHandler(srv)
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				log.Fatal("%w", err)
			}
		}()
	default:
		fmt.Println("The wrong number!")
		defer os.Exit(1)
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
