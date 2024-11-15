package main

import (
	"fmt"
	"net/http"

	"github.com/ksel172/Meduza/teamserver/internal/server"
)

func main() {

	// NewServer initialize the Http Server
	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
