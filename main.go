// Package main contains main function
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/distuurbia/firstTask/internal/config"
	"github.com/distuurbia/firstTask/internal/handler"
	customMidleware "github.com/distuurbia/firstTask/internal/middleware"
	"github.com/distuurbia/firstTask/internal/repository"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/caarlos0/env/v8"
)

// ConnectPgx connects to the pgxpool
func ConnectPgx() (*pgxpool.Pool, error) {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	cfgPgx, err := pgxpool.ParseConfig(cfg.PgxConnectionString)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfgPgx)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

// ConnectMongo connects to the mongoDB
func ConnectMongo() (*mongo.Client, error) {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	const ctxTimeout = 10
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConnectionString))
	if err != nil {
		return client, fmt.Errorf("%w", err)
	}
	return client, nil
}

// main is an executable function
func main() {
	var handl *handler.EntityHandler
	validate := validator.New()
	fmt.Println("What db do u wanna use?\n 1.PostgreSQL\n 2.MongoDB")
	var dbChoose int
	_, err := fmt.Scan(&dbChoose)
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
		persPgx := repository.NewRepositoryPgx(dbpool)
		persSrv := service.NewPersonService(persPgx)
		userSrv := service.NewUserService(persPgx)
		handl = handler.NewHandler(persSrv, userSrv, validate)
	case MongoDB:
		client, err := ConnectMongo()
		if err != nil {
			log.Fatal("could not construct the client: ", err)
		}
		rpsMongo := repository.NewRepositoryMongo(client)
		srvPers := service.NewPersonService(rpsMongo)
		srvUser := service.NewUserService(rpsMongo)
		handl = handler.NewHandler(srvPers, srvUser, validate)
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				log.Fatal("could not disconnect ", err)
			}
		}()
	default:
		//nolint:gocritic
		log.Fatal("The wrong number!")
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/persondb", handl.Create, customMidleware.JWTMiddleware())
	e.GET("/persondb/:id", handl.ReadRow, customMidleware.JWTMiddleware())
	e.PUT("/persondb/:id", handl.Update, customMidleware.JWTMiddleware())
	e.DELETE("/persondb/:id", handl.Delete, customMidleware.JWTMiddleware())

	e.POST("/signUp", handl.SignUp)
	e.POST("/login", handl.Login)
	e.POST("/refresh", handl.Refresh)
	e.Logger.Fatal(e.Start(":8080"))
}
