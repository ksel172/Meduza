package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/internal/database"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type Server struct {
	port int

	db database.Service
}

func NewServer() *http.Server {
	port := utils.GetEnvInt("PORT", 8080)
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
