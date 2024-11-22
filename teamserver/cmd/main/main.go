package main

import (
	"github.com/ksel172/Meduza/teamserver/internal/models"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"log"
	"os"
	"time"

	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/server"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
)

func main() {

	// Initialize services
	log.Println("Connecting to postgres db...")
	pgsql, err := repos.Setup()
	if err != nil {
		log.Fatal("Failed to connect to pgsql. Terminating...", err)
	}
	defer pgsql.Close()
	log.Println("Connected to postgres db")

	log.Println("Connecting to redisService db...")
	redisService := repos.NewRedisService()
	log.Println("Connected to redisService db")

	log.Println("Setting up data access layers...")
	userDal := dal.NewUsersDAL(pgsql, conf.GetMeduzaDbSchema())
	adminDal := dal.NewAdminsDAL(pgsql, conf.GetMeduzaDbSchema())
	agentDal := dal.NewAgentDAL(&redisService)
	checkInDal := dal.NewCheckInDAL(&redisService)
	log.Println("Finished setting up data access layers")

	secret := os.Getenv("JWT_SECRET")
	jwtService := models.NewJWTService(secret, 15*time.Minute, 30*24*time.Hour)

	userController := handlers.NewUserController(userDal)
	authController := handlers.NewAuthController(userDal, jwtService)
	adminController := handlers.NewAdminController(adminDal)
	agentController := handlers.NewAgentController(agentDal)
	checkInController := handlers.NewCheckInController(checkInDal, agentDal)

	dependencies := &server.DependencyContainer{
		UserController:    userController,
		RedisService:      &redisService,
		AuthController:    authController,
		JwtService:        jwtService,
		AdminController:   adminController,
		AgentController:   agentController,
		CheckInController: checkInController,
	}

	// NewServer initialize the Http Server
	teamserver := server.NewServer(dependencies)

	log.Println("Starting teamserver...")
	if err := teamserver.Run(); err != nil {
		log.Panic("Failed to start teamserver. Terminating...", err)
	}
}
