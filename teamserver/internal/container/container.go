package container

import (
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"

	// services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

// Controllers owned by the listener server
type Container struct {
	UserController     *handlers.UserController
	RedisService       *repos.Service
	AuthController     *handlers.AuthController
	TeamController     *handlers.TeamController
	JwtService         models.JWTServiceProvider
	AgentController    *handlers.AgentController
	ListenerController *services.ListenerController
	// ListenerService       *services.ListenersService // for autostart
	// ListenerDal           *dal.ListenerDAL
	PayloadController     *handlers.PayloadHandler
	ModuleController      *handlers.ModuleController
	CertificateController *handlers.CertificateHandler
}

func NewContainer() (*Container, error) {
	logger.Info("Connecting to Postgres db...")
	pgsql, err := repos.Setup()
	if err != nil {
		logger.Error("Error while setting Up Postgres:", err)
		return nil, err
	}

	logger.Info("Setting Up Data Access logLevel")
	schema := conf.GetMeduzaDbSchema()
	userDal := dal.NewUsersDAL(pgsql, schema)
	teamDal := dal.NewTeamDAL(pgsql, schema)
	agentDal := dal.NewAgentDAL(pgsql, schema)
	// checkInDal := dal.NewCheckInDAL(pgsql, schema)
	listenerDal := services.NewListenerDAL(pgsql, schema)
	listenerManager := services.NewListenerManager(listenerDal)
	payloadDal := dal.NewPayloadDAL(pgsql, schema)
	moduleDal := dal.NewModuleDAL(pgsql, schema)
	certificateDal := dal.NewCertificateDAL(pgsql, schema)

	// Initialize services
	redisService := repos.NewRedisService()
	jwtService := models.NewJWTService(conf.GetMeduzaJWTToken(), 30*time.Minute, 30*24*time.Hour)

	//Type assertion error fix
	// autoStart, ok := listenerDal.(*dal.ListenerDAL)
	// if !ok {
	// 	logger.Warn("Unable to type assetion ListenerDAL")
	// }

	return &Container{
		UserController:     handlers.NewUserController(userDal),
		RedisService:       &redisService,
		AuthController:     handlers.NewAuthController(userDal, jwtService),
		TeamController:     handlers.NewTeamController(teamDal),
		JwtService:         jwtService,
		AgentController:    handlers.NewAgentController(agentDal, moduleDal),
		ListenerController: services.NewListenerController(listenerManager),
		// ListenerController:    handlers.NewListenersHandler(listenerDal, listenersService),
		// ListenerService:       listenersService,
		// ListenerDal:           autoStart,
		PayloadController:     handlers.NewPayloadHandler(agentDal, listenerDal, payloadDal),
		ModuleController:      handlers.NewModuleController(moduleDal),
		CertificateController: handlers.NewCertificateHandler(certificateDal),
		// ListenerContainer: ListenerContainer{
		// 	CheckInController: checkInController,
		// },
	}, nil
}
