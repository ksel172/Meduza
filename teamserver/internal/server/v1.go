package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
)

func (s *Server) AuthV1(group *gin.RouterGroup) {
	authRoutes := group.Group("/auth")
	{
		authRoutes.POST("/register", s.dependencies.UserController.AddUsers)
		authRoutes.POST("/add-admin", s.dependencies.AdminController.CreateAdmin)
		authRoutes.POST("/login", s.dependencies.AuthController.LoginController)
		authRoutes.GET("/refresh-token", s.dependencies.AuthController.RefreshTokenController)
		authRoutes.POST("/logout", s.dependencies.AuthController.LogoutController)
	}
}

func (s *Server) AdminV1(group *gin.RouterGroup) {

	adminProtectedRoutes := group.Group("/teamserver")
	{
		adminProtectedRoutes.Use(s.AdminMiddleware())
		adminProtectedRoutes.GET("/users", s.dependencies.UserController.GetUsers)
		adminProtectedRoutes.POST("/users", s.dependencies.UserController.AddUsers)
	}
}

func (s *Server) AgentsV1(group *gin.RouterGroup) {

	agentsGroup := group.Group("/agents")
	{
		agentsGroup.GET(fmt.Sprintf("/:%s", models.ParamAgentID), s.dependencies.AgentController.GetAgent)
		agentsGroup.PUT(fmt.Sprintf("/:%s", models.ParamAgentID), s.dependencies.AgentController.UpdateAgent)
		agentsGroup.DELETE(fmt.Sprintf("/:%s", models.ParamAgentID), s.dependencies.AgentController.DeleteAgent)

		// Agent Tasks API
		agentsGroup.GET("/tasks", s.dependencies.AgentController.GetAgentTasks)
		agentsGroup.POST("/tasks", s.dependencies.AgentController.CreateAgentTask)
		agentsGroup.DELETE(fmt.Sprintf(":%s/tasks", models.ParamAgentID), s.dependencies.AgentController.DeleteAgentTasks)
		agentsGroup.DELETE(fmt.Sprintf(":%s/tasks/:%s", models.ParamAgentID, models.ParamTaskID), s.dependencies.AgentController.DeleteAgentTask)
	}
}

func (s *Server) CheckInV1(group *gin.RouterGroup) {

	checkinGroup := group.Group("/checkin")
	{
		checkinGroup.POST("/", s.dependencies.CheckInController.CreateAgent)
		checkinGroup.GET("/", s.dependencies.CheckInController.GetTasks)
	}
}

func (s *Server) ListenersV1(group *gin.RouterGroup) {

	listenersGroup := group.Group("/listeners")
	{
		listenersGroup.POST("", s.dependencies.ListenerController.AddListener)
		listenersGroup.GET("/:id", s.dependencies.ListenersController.GetListener)
		listenersGroup.PUT("", s.dependencies.ListenersController.UpdateListener)
		listenersGroup.DELETE(":id", s.dependencies.ListenersController.DeleteListener)
	}
}
