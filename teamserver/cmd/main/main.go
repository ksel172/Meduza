package main

import (
	"github.com/ksel172/Meduza/teamserver/internal/container"
	"github.com/ksel172/Meduza/teamserver/internal/server"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

func main() {
	dependencies, err := container.NewContainer()
	if err != nil {
		logger.Error("Error While Setting dependencies", err)
		return
	}

	// NewServer initialize the Http Server
	teamserver := server.NewServer(dependencies)

	teamserver.RegisterRoutes()

	logger.Info("Starting Teamserver...")
	if err := teamserver.Run(); err != nil {
		logger.Panic("Failed to Start Teamserver. Terminating...", err)
	}

}
