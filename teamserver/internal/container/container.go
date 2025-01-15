package container

import (
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

// Controllers owned by the listener server
type ListenerContainer struct {
	CheckInController   *services.CheckInController
	AgentAuthController *services.AgentAuthController
}

type Container struct {
	UserController     *handlers.UserController
	RedisService       *repos.Service
	AuthController     *handlers.AuthController
	JwtService         models.JWTServiceProvider
	AdminController    *handlers.AdminController
	AgentController    *handlers.AgentController
	ListenerController *handlers.ListenerHandler
	PayloadController  *handlers.PayloadHandler
	ListenerContainer
}

func NewContainer() (*Container, error) {
	logger.Info("Connecting to Postgres db...")
	pgsql, err := repos.Setup()
	if err != nil {
		logger.Error("Error while setting Up Postgres:", err)
		return nil, err
	}

	logger.Info("Setting Up Data Access Layer")
	schema := conf.GetMeduzaDbSchema()
	userDal := dal.NewUsersDAL(pgsql, schema)
	adminDal := dal.NewAdminsDAL(pgsql, schema)
	agentDal := dal.NewAgentDAL(pgsql, schema)
	checkInDal := dal.NewCheckInDAL(pgsql, schema)
	listenerDal := dal.NewListenerDAL(pgsql, schema)
	payloadDal := dal.NewPayloadDAL(pgsql, schema)

	// CheckInController is owned by the listener server
	checkInController := services.NewCheckInController(checkInDal, agentDal)
	agentAuthController := services.NewAgentAuthController(payloadDal)

	// Initialize services
	redisService := repos.NewRedisService()
	jwtService := models.NewJWTService(conf.GetMeduzaJWTToken(), 15, 30*24*60*60)
	listenersService := services.NewListenerService(checkInController, agentAuthController)

	return &Container{
		UserController:     handlers.NewUserController(userDal),
		RedisService:       &redisService,
		AuthController:     handlers.NewAuthController(userDal, jwtService),
		JwtService:         jwtService,
		AdminController:    handlers.NewAdminController(adminDal),
		AgentController:    handlers.NewAgentController(agentDal),
		ListenerController: handlers.NewListenersHandler(listenerDal, listenersService),
		PayloadController:  handlers.NewPayloadHandler(agentDal, listenerDal, payloadDal),
		ListenerContainer: ListenerContainer{
			CheckInController:   checkInController,
			AgentAuthController: agentAuthController,
		},
	}, nil
}
