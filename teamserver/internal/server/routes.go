package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/services/api"
)

func (s *Server) RegisterRoutes() http.Handler {

	router := gin.Default()

	router.GET("/", api.HelloWorldHandler)

	// base URL
	apiGroup := router.Group("/api")
	checkinGroup := router.Group("/checkin")

	// Version control
	v1Group := apiGroup.Group("/v1")

	// Groups
	usersGroup := v1Group.Group("/users")
	agentsGroup := v1Group.Group("/agents")

	// Users API
	usersGroup.GET("/", func(context *gin.Context) {
		s.dependencies.UserController.GetUsers(context.Writer, context.Request)
	})

	// Agents API
	agentsGroup.GET("/", func(context *gin.Context) {
		s.dependencies.AgentController.GetAgent(context.Writer, context.Request)
	})
	agentsGroup.PUT("/", func(context *gin.Context) {
		s.dependencies.AgentController.UpdateAgent(context.Writer, context.Request)
	})
	agentsGroup.DELETE("/", func(context *gin.Context) {
		s.dependencies.AgentController.DeleteAgent(context.Writer, context.Request)
	})

	// AgentTasks API
	agentsGroup.GET("/tasks", func(ctx *gin.Context) {
		s.dependencies.AgentController.GetAgentTasks(ctx.Writer, ctx.Request)
	})
	agentsGroup.POST("/tasks", func(context *gin.Context) {
		s.dependencies.AgentController.CreateAgentTask(context.Writer, context.Request)
	})
	agentsGroup.DELETE("/tasks", func(context *gin.Context) {
		s.dependencies.AgentController.DeleteAgentTasks(context.Writer, context.Request)
	})
	agentsGroup.DELETE("/tasks/task", func(context *gin.Context) {
		s.dependencies.AgentController.DeleteAgentTask(context.Writer, context.Request)
	})

	// CheckIn Controller API
	checkinGroup.POST("/", func(context *gin.Context) {
		s.dependencies.CheckInController.CreateAgent(context.Writer, context.Request)
	})
	checkinGroup.GET("/", func(context *gin.Context) {
		s.dependencies.CheckInController.GetTasks(context.Writer, context.Request)
	})

	return router
}
