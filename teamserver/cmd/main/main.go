package main

import (
	"database/sql"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/app/users"
	"github.com/ksel172/Meduza/teamserver/internal/server"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"log"
)

func main() {

	// Initialize services
	log.Println("Connecting to postgres db...")
	pgsql, err := storage.Setup()
	if err != nil {
		log.Fatal("Failed to connect to pgsql. Terminating...", err)
	}
	defer pgsql.Close()
	log.Println("Connected to postgres db")

	log.Println("Connecting to redisService db...")
	redisService := redis.NewRedisService()
	log.Println("Connected to redisService db")

	log.Println("Loading users service...")
	userDal := dal.NewUsersDAL(pgsql, conf.GetMeduzaDbSchema())
	userService := users.NewService(userDal)
	userController := handlers.NewUserController(userService)

	// Create dependency container
	dependencies := &server.DependencyContainer{
		UserController: userController,
		UserService:    userService,
		RedisService:   redisService,
	}

	// NewServer initialize the Http Server
	teamserver := server.NewServer(dependencies)

	log.Println("Starting teamserver...")
	if err := teamserver.Run(); err != nil {
		log.Panicf("Failed to start teamserver. Terminating...", err)
	}
}
