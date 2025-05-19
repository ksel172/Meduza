package container

import (
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	listener_service "github.com/ksel172/Meduza/teamserver/internal/services/listener"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"

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
	ListenerController *handlers.ListenerController
	// ListenerService       *services.ListenersService // for autostart
	// ListenerDal           *dal.ListenerDAL
	PayloadController          *handlers.PayloadController
	ModuleController           *handlers.ModuleController
	CertificateController      *handlers.CertificateHandler
	ExternalListenerController *handlers.ExternalController
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
	listenerDal := dal.NewListenerDAL(pgsql, schema)
	moduleDal := dal.NewModuleDAL(pgsql, schema)
	certificateDal := dal.NewCertificateDAL(pgsql, schema)
	payloadDal := dal.NewPayloadDAL(pgsql, schema)
	// Initialize services
	redisService := repos.NewRedisService()
	jwtService := models.NewJWTService(conf.GetMeduzaJWTToken(), 30*time.Minute, 30*24*time.Hour)
	listenerService := listener_service.NewListenerService(listenerDal)
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
		ListenerController: handlers.NewListenersHandler(listenerService, listenerDal),
		// ListenerController:    handlers.NewListenersHandler(listenerDal, listenersService),
		// ListenerService:       listenersService,
		// ListenerDal:           autoStart,
		PayloadController:     handlers.NewPayloadController(payloadDal),
		ModuleController:      handlers.NewModuleController(moduleDal),
		CertificateController: handlers.NewCertificateHandler(certificateDal),
		// ListenerContainer: ListenerContainer{
		// 	CheckInController: checkInController,
		// },
		ExternalListenerController: handlers.NewExternalController(agentDal, listenerDal, checkin.CheckInController{}),
	}, nil
}
