// Package main contains main function
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/distuurbia/firstTask/docs"

	"github.com/caarlos0/env/v8"
	"github.com/distuurbia/firstTask/internal/config"
	"github.com/distuurbia/firstTask/internal/handler"
	customMidleware "github.com/distuurbia/firstTask/internal/middleware"
	"github.com/distuurbia/firstTask/internal/repository"
	"github.com/distuurbia/firstTask/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectPgx connects to the pgxpool
func ConnectPgx(cfg *config.Config) (*pgxpool.Pool, error) {
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
func ConnectMongo(cfg *config.Config) (*mongo.Client, error) {
	const ctxTimeout = 10
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoConnectionString))
	if err != nil {
		return client, fmt.Errorf("%w", err)
	}
	return client, nil
}

// ConnectRedis connects to the redis db
func ConnectRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	return client
}

// @title FirstTask API
// @description API for managing persons and users
// @version 1.0
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		//nolint:gocritic
		log.Fatal("could not parse config: ", err)
	}
	var handl *handler.EntityHandler
	validate := validator.New()
	rdsClient := ConnectRedis(&cfg)
	rds := repository.NewRepositoryRedis(rdsClient)
	fmt.Println("What db do u wanna use?\n 1.PostgreSQL\n 2.MongoDB")
	var dbChoose int
	// _, err := fmt.Scan(&dbChoose)
	// if err != nil {
	// 	fmt.Println("failed to scan")
	// }
	dbChoose = 1
	const PostgreSQL = 1
	const MongoDB = 2
	switch dbChoose {
	case PostgreSQL:
		dbpool, err := ConnectPgx(&cfg)
		if err != nil {
			log.Fatal("could not construct the pool: ", err)
		}
		defer dbpool.Close()
		persPgx := repository.NewRepositoryPgx(dbpool)
		persSrv := service.NewPersonService(persPgx, rds)
		userSrv := service.NewUserService(persPgx, &cfg)
		handl = handler.NewHandler(persSrv, userSrv, validate)
	case MongoDB:
		client, err := ConnectMongo(&cfg)
		if err != nil {
			//nolint:gocritic
			log.Fatal("could not construct the client: ", err)
		}
		rpsMongo := repository.NewRepositoryMongo(client)
		srvPers := service.NewPersonService(rpsMongo, rds)
		srvUser := service.NewUserService(rpsMongo, &cfg)
		handl = handler.NewHandler(srvPers, srvUser, validate)
		defer func() {
			if err = client.Disconnect(context.Background()); err != nil {
				//nolint:gocritic
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

	e.POST("/persons", handl.Create, customMidleware.JWTMiddleware(&cfg))
	e.GET("/persons/:id", handl.ReadRow, customMidleware.JWTMiddleware(&cfg))
	e.GET("/persons", handl.GetAll, customMidleware.JWTMiddleware(&cfg))
	e.PUT("/persons/:id", handl.Update, customMidleware.JWTMiddleware(&cfg))
	e.DELETE("/persons/:id", handl.Delete, customMidleware.JWTMiddleware(&cfg))

	e.POST("/signUp", handl.SignUp)
	e.POST("/login", handl.Login)
	e.POST("/refresh", handl.Refresh)
	e.GET("/downloadImage/:imageName", handl.DownloadImage)
	e.POST("/uploadImage", handl.UploadImage)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
