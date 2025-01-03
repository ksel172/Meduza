package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/internal/container"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/mattn/go-colorable"
)

type Server struct {
	host         string
	port         int
	engine       *gin.Engine
	dependencies *container.Container
}

func NewServer(dependencies *container.Container) *Server {

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

	return server
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.engine.Run(addr)
}
