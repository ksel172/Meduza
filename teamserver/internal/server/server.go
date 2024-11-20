package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/app/users"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
)

type DependencyContainer struct {
	UserController *handlers.UserController
	UserService    *users.Service
	RedisService   *redis.Service
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
