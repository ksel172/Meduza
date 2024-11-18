package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/services/api"
)

func (s *Server) RegisterRoutes() http.Handler {

	router := gin.Default()

	router.GET("/", api.HelloWorldHandler)

	router.POST("/add-users", s.dependencies.UserController.AddUsersController)

	router.POST("/login", s.dependencies.AuthController.LoginController)

	protected := router.Group("/api/v1/teamserver")

	protected.Use(s.HandleCors())
	protected.Use(s.JWTAuthMiddleware())

	protected.GET("/users", s.dependencies.UserController.GetUsersController)

	return router
}
