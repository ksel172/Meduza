package server

import (
	"fmt"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/services/api"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type DependencyContainer struct {
	UserController *api.UserController
}

type Server struct {
	host        string
	port        int
	depenencies *DependencyContainer
}

func NewServer(dependencies *DependencyContainer) *http.Server {
	NewServer := &Server{
		host:        conf.GetMeduzaServerHostname(),
		port:        conf.GetMeduzaServerPort(),
		depenencies: dependencies,
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
