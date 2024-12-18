package main

import (
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"

	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/server"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
)

func main() {
	logger.Info("Connecting to Postgres db...")

	pgsql, err := repos.Setup()
	if err != nil {
		logger.Fatal("Failed to connect to Pgsql. Terminating...", err)
	}
	defer pgsql.Close()
	logger.Good("Connected to Postgres DB")

	logger.Info("Connecting to RedisService Db...")
	redisService := repos.NewRedisService()
	logger.Good("Connected to RedisService Db")

	logger.Info("Setting up Data Access Layers...")
	userDal := dal.NewUsersDAL(pgsql, conf.GetMeduzaDbSchema())
	adminDal := dal.NewAdminsDAL(pgsql, conf.GetMeduzaDbSchema())
	agentDal := dal.NewAgentDAL(&redisService)
	checkInDal := dal.NewCheckInDAL(&redisService)
	listenerDal := dal.NewListenersDAL(pgsql, conf.GetMeduzaDbSchema())
	logger.Good("Finished setting up data access layers")

	secret := conf.GetMeduzaJWTToken()
	jwtService := models.NewJWTService(secret, 15*time.Minute, 30*24*time.Hour)

	userController := handlers.NewUserController(userDal)
	authController := handlers.NewAuthController(userDal, jwtService)
	adminController := handlers.NewAdminController(adminDal)
	agentController := handlers.NewAgentController(agentDal)
	checkInController := handlers.NewCheckInController(checkInDal, agentDal)
	listenerController := handlers.NewListenersHandler(listenerDal)

	dependencies := &server.DependencyContainer{
		UserController:     userController,
		RedisService:       &redisService,
		AuthController:     authController,
		JwtService:         jwtService,
		AdminController:    adminController,
		AgentController:    agentController,
		CheckInController:  checkInController,
		ListenerController: listenerController,
	}

	// NewServer initialize the Http Server
	teamserver := server.NewServer(dependencies)

	logger.Info("Starting Teamserver...")
	if err := teamserver.Run(); err != nil {
		logger.Panic("Failed to Start Teamserver. Terminating...", err)
	}
}
