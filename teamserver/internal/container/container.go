package container

import (
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/services"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type Container struct {
	UserController     *handlers.UserController
	RedisService       *repos.Service
	AuthController     *handlers.AuthController
	JwtService         models.JWTServiceProvider
	AdminController    *handlers.AdminController
	AgentController    *handlers.AgentController
	CheckInController  *handlers.CheckInController
	ListenerController *handlers.ListenerHandler
}

func NewContainer() (*Container, error) {
	logger.Info("Connecting to Postgres db...")
	pgsql, err := repos.Setup()
	if err != nil {
		logger.Error("Error while setting Up Postgres:", err)
		return nil, err
	}

	redisService := repos.NewRedisService()
	schema := conf.GetMeduzaDbSchema()

	logger.Info("Setting Up Data Access Layer")
	userDal := dal.NewUsersDAL(pgsql, schema)
	adminDal := dal.NewAdminsDAL(pgsql, schema)
	agentDal := dal.NewAgentDAL(pgsql, schema)
	checkInDal := dal.NewCheckInDAL(pgsql, schema)
	listenerDal := dal.NewListenerDAL(pgsql, schema)

	jwtService := models.NewJWTService(conf.GetMeduzaJWTToken(), 15, 30*24*60*60)
	listenersService := services.NewListenerService()

	return &Container{
		UserController:     handlers.NewUserController(userDal),
		RedisService:       &redisService,
		AuthController:     handlers.NewAuthController(userDal, jwtService),
		JwtService:         jwtService,
		AdminController:    handlers.NewAdminController(adminDal),
		AgentController:    handlers.NewAgentController(agentDal),
		CheckInController:  handlers.NewCheckInController(checkInDal, agentDal),
		ListenerController: handlers.NewListenersHandler(listenerDal, listenersService),
	}, nil
}
