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
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/services/api"
	"github.com/ksel172/Meduza/teamserver/services/auth"
	"github.com/ksel172/Meduza/teamserver/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could'nt load .env file")
	}

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

	secret := os.Getenv("JWT_SECRET")

	userDal := storage.NewUsersDAL(postgres, conf.GetMeduzaDbSchema())
	adminDal := storage.NewAdminsDAL(postgres, conf.GetMeduzaDbSchema())
	userController := api.NewUserController(userDal)
	jwtService := auth.NewJWTService(secret, 15*time.Minute, 30*24*time.Hour)
	authController := api.NewAuthController(userDal, jwtService)
	adminController := api.NewAdminController(adminDal)

	agentDal := redis.NewAgentDAL(redisService)
	agentController := api.NewAgentController(agentDal)

	checkInDal := redis.NewCheckInDAL(redisService)
	checkInController := api.NewCheckInController(checkInDal)

	return &server.DependencyContainer{
		UserController:  userController,
		RedisService:    redis,
		AuthController:  authController,
		JwtService:      jwtService,
		AdminController: adminController,
		AgentController:   agentController,
		CheckInController: checkInController,
	}
}
