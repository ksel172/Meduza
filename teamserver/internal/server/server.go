package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/storage/redis"
	"github.com/ksel172/Meduza/teamserver/services/api"

	_ "github.com/joho/godotenv/autoload"
)

type DependencyContainer struct {
	UserController    *api.UserController
	AgentController   *api.AgentController
	CheckInController *api.CheckInController
	RedisService      *redis.Service
}

type Server struct {
	host         string
	port         int
	dependencies *DependencyContainer
}

func NewServer(dependencies *DependencyContainer) *http.Server {
	NewServer := &Server{
		host:         conf.GetMeduzaServerHostname(),
		port:         conf.GetMeduzaServerPort(),
		dependencies: dependencies,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", NewServer.host, NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
