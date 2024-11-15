package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/services/api"
)

func (s *Server) RegisterRoutes() http.Handler {

	router := gin.Default()

	router.GET("/", api.HelloWorldHandler)
	router.GET("/api/users", api.UserController{})

	return router
}
