package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/app/users"
	"github.com/ksel172/Meduza/teamserver/internal/server"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/services/api"
	"github.com/ksel172/Meduza/teamserver/services/auth"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could'nt load .env file")
	}

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
	secret := os.Getenv("JWT_SECRET")
	adminDal := storage.NewAdminsDAL(pgsql, conf.GetMeduzaDbSchema())
	jwtService := auth.NewJWTService(secret, 15*time.Minute, 30*24*time.Hour)
	authController := api.NewAuthController(userDal, jwtService)
	adminController := api.NewAdminController(adminDal)
	agentDal := redis.NewAgentDAL(redisService)
	checkInDal := redis.NewCheckInDAL(redisService)
	agentController := api.NewAgentController(agentDal)
	checkInController := api.NewCheckInController(checkInDal, agentDal)

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
