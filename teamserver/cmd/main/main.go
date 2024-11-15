package main

import (
	"fmt"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/services/api"
	"log"
	"net/http"

	"github.com/ksel172/Meduza/teamserver/internal/server"
)

func main() {

	// Initialize services
	log.Println("Connecting to database...")
	database, err := storage.Setup()
	if err != nil {
		log.Fatal("Failed to connect to database. Terminating...", err)
	}
	defer database.Close()
	log.Println("Connected to database")

	userDal := storage.NewUsersDAL(database, conf.GetMeduzaDbSchema())

	userController := api.NewUserController(userDal)

	// NewServer initialize the Http Server
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
