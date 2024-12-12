package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/mattn/go-colorable"
)

type DependencyContainer struct {
	UserController      *handlers.UserController
	RedisService        *repos.Service
	AuthController      *handlers.AuthController
	JwtService          *models.JWTService
	AdminController     *handlers.AdminController
	AgentController     *handlers.AgentController
	CheckInController   *handlers.CheckInController
	ListenersController *handlers.ListenerHandler
}

type Server struct {
	host         string
	port         int
	engine       *gin.Engine
	dependencies *DependencyContainer
}

func NewServer(dependencies *DependencyContainer) *Server {

	// it force the console output to be colored.
	gin.ForceConsoleColor()

	// Declare Server config
	server := &Server{
		host:         conf.GetMeduzaServerHostname(),
		port:         conf.GetMeduzaServerPort(),
		engine:       gin.Default(),
		dependencies: dependencies,
	}

	gin.DefaultWriter = colorable.NewColorableStdout()

	server.RegisterRoutes()

	return server
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.engine.Run(addr)
}
