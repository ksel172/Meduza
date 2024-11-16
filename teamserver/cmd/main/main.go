package main

import (
	"database/sql"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/services/api"
	"log"
	"net/http"

	"github.com/ksel172/Meduza/teamserver/internal/server"
)

func main() {

	// Initialize services
	log.Println("Connecting to postgres db...")
	database, err := storage.Setup()
	if err != nil {
		log.Fatal("Failed to connect to database. Terminating...", err)
	}
	defer database.Close()
	log.Println("Connected to postgres db")

	log.Println("Connecting to redisService db...")
	redisService := redis.NewRedisService()
	log.Println("Connected to redisService db")

	// Create dependency container
	dependencies := InitializeDependencies(database, &redisService)

	// NewServer initialize the Http Server
	newServer := server.NewServer(dependencies)

	err = newServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http newServer error: %s", err))
	}
}

func InitializeDependencies(postgres *sql.DB, redis *redis.Service) *server.DependencyContainer {
	userDal := storage.NewUsersDAL(postgres, conf.GetMeduzaDbSchema())
	userController := api.NewUserController(userDal)

	return &server.DependencyContainer{
		UserController: userController,
		RedisService:   redis,
	}
}
