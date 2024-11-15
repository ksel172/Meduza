package server

import (
	"fmt"
	"github.com/ksel172/Meduza/teamserver/conf"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
)

type Server struct {
	host string
	port int
	db   storage.Service
}

func NewServer() *http.Server {
	NewServer := &Server{
		host: conf.GetMeduzaServerHostname(),
		port: conf.GetMeduzaServerPort(),
		db:   storage.New(),
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
