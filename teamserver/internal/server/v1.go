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
		agentsGroup.POST("/tasks/:id", s.dependencies.AgentController.CreateAgentTask)
		agentsGroup.DELETE(fmt.Sprintf(":%s/tasks", models.ParamAgentID), s.dependencies.AgentController.DeleteAgentTasks)
		agentsGroup.DELETE(fmt.Sprintf(":%s/tasks/:%s", models.ParamAgentID, models.ParamTaskID), s.dependencies.AgentController.DeleteAgentTask)

		// Agent config API
		agentsGroup.POST("/config/:id", s.dependencies.AgentController.CreateAgentConfig)
		agentsGroup.PUT("/config/:id", s.dependencies.AgentController.UpdateAgentConfig)
		agentsGroup.GET("/config/:id", s.dependencies.AgentController.GetAgentConfig)
		agentsGroup.DELETE("/config/:id", s.dependencies.AgentController.DeleteAgentConfig)

		// Agent info API
		agentsGroup.POST("/info", s.dependencies.AgentController.CreateAgentInfo)
		agentsGroup.PUT("/info/:id", s.dependencies.AgentController.UpdateAgentInfo)
		agentsGroup.GET("/info/:id", s.dependencies.AgentController.GetAgentInfo)
		agentsGroup.DELETE("/info/:id", s.dependencies.AgentController.DeleteAgentInfo)
	}
}
func (s *Server) ListenersV1(group *gin.RouterGroup) {

	listenersGroup := group.Group("/listeners")
	{
		listenersGroup.POST("", s.dependencies.ListenerController.CreateListener) // pg
		listenersGroup.GET("/:id", s.dependencies.ListenerController.GetListenerById)
		listenersGroup.GET("/all", s.dependencies.ListenerController.GetAllListeners)
		listenersGroup.PUT("/:id", s.dependencies.ListenerController.UpdateListener)
		listenersGroup.DELETE("/:id", s.dependencies.ListenerController.DeleteListener)
		listenersGroup.POST("/:id/start", s.dependencies.ListenerController.StartListener)
		listenersGroup.POST("/:id/stop", s.dependencies.ListenerController.StopListener)
		listenersGroup.GET("/:id/status", s.dependencies.ListenerController.CheckRunningListener)
	}
}
func (s *Server) PayloadV1(group *gin.RouterGroup) {

	payloadsGroup := group.Group("/payloads")
	{
		payloadsGroup.POST("/create", s.dependencies.PayloadController.CreatePayload)
		payloadsGroup.GET("/all", s.dependencies.PayloadController.GetAllPayloads)
		payloadsGroup.POST("/delete/:id", s.dependencies.PayloadController.DeletePayload)
		payloadsGroup.GET("/download/:id", s.dependencies.PayloadController.DownloadPayload)
		payloadsGroup.POST("/delete/all", s.dependencies.PayloadController.DeleteAllPayloads)
	}
}

func (s *Server) ModuleV1(group *gin.RouterGroup) {

	moduleGroup := group.Group("/modules")
	{
		moduleGroup.POST("/upload", s.dependencies.ModuleController.UploadModule)
		moduleGroup.POST("/delete/:id", s.dependencies.ModuleController.DeleteModule)
		moduleGroup.POST("/delete/all", s.dependencies.ModuleController.DeleteAllModules)
		moduleGroup.GET("/all", s.dependencies.ModuleController.GetAllModules)
		moduleGroup.GET("/:id", s.dependencies.ModuleController.GetModuleById)
	}
}
