package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/app/users"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/services/api"
	"github.com/ksel172/Meduza/teamserver/services/auth"

	_ "github.com/joho/godotenv/autoload"
)

type DependencyContainer struct {
	UserController    *api.UserController
	RedisService      *redis.Service
	AuthController    *api.AuthController
	JwtService        *auth.JWTService
	AdminController   *api.AdminController
	AgentController   *api.AgentController
	CheckInController *api.CheckInController
	UserController    *handlers.UserController
	UserService       *users.Service
	RedisService      *redis.Service
}

type Server struct {
	host         string
	port         int
	engine       *gin.Engine
	dependencies *DependencyContainer
}

func NewServer(dependencies *DependencyContainer) *Server {

	// Declare Server config
	server := &Server{
		host:         fmt.Sprintf("%s:%d", conf.GetMeduzaServerHostname(), conf.GetMeduzaServerPort()),
		engine:       gin.Default(),
		dependencies: dependencies,
	}

	server.RegisterRoutes()

	return server
}

func (s *Server) RegisterRoutes() {
	users := s.engine.Group("api/v1/users")
	{
		users.GET("", s.dependencies.UserController.GetUsers)
	}

	// TODO add listeners
	//listeners := s.engine.Group("api/v1/listeners")
	//{
	//	listeners.GET("", s.dependencies.)
	//}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.engine.Run(addr)
}
