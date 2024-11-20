package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {

	router := gin.Default()

	pubRoutes := router.Group("/api/v1/teamserver/public") // public routes

	{
		pubRoutes.Use(s.HandleCors())
		pubRoutes.POST("/register", s.dependencies.UserController.AddUsersController)
		pubRoutes.POST("/add-admin", s.dependencies.AdminController.CreateAdmin)
	}

	authroutes := router.Group("/api/v1/auth") // auth routes

	{
		pubRoutes.Use(s.HandleCors())
		authroutes.POST("/sign-in", s.dependencies.AuthController.LoginController)
		authroutes.GET("/refresh-token", s.dependencies.AuthController.RefreshTokenController)
	}

	adminProtected := router.Group("/api/v1/teamserver") // admin protected routes
	{
		adminProtected.Use(s.HandleCors())
		adminProtected.Use(s.AdminMiddleware())
		adminProtected.GET("/users", s.dependencies.UserController.GetUsersController)
	}

	return router
}
