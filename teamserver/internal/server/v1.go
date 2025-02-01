package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
)

func (s *Server) AuthV1(group *gin.RouterGroup) {
	authRoutes := group.Group("/auth")
	{
		authRoutes.POST("/register", s.AdminMiddleware(), s.dependencies.UserController.AddUsers)
		authRoutes.POST("/login", s.dependencies.AuthController.LoginController)
		authRoutes.GET("/refresh", s.UserMiddleware(), s.dependencies.AuthController.RefreshTokenController)
		authRoutes.POST("/logout", s.UserMiddleware(), s.dependencies.AuthController.LogoutController)
	}
}

func (s *Server) UsersV1(group *gin.RouterGroup) {

	adminProtectedRoutes := group.Group("/users")
	{
		adminProtectedRoutes.Use(s.AdminMiddleware())
		adminProtectedRoutes.GET("", s.dependencies.UserController.GetUsers)
		adminProtectedRoutes.POST("", s.dependencies.UserController.AddUsers)
	}
}

func (s *Server) AgentsV1(group *gin.RouterGroup) {

	agentsGroup := group.Group("/agents")
	{
		agentsGroup.Use(s.UserMiddleware())
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
		listenersGroup.Use(s.UserMiddleware())
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
		payloadsGroup.Use(s.UserMiddleware())
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
		moduleGroup.Use(s.UserMiddleware())
		moduleGroup.POST("/upload", s.dependencies.ModuleController.UploadModule)
		moduleGroup.POST("/delete/:id", s.dependencies.ModuleController.DeleteModule)
		moduleGroup.POST("/delete/all", s.dependencies.ModuleController.DeleteAllModules)
		moduleGroup.GET("/all", s.dependencies.ModuleController.GetAllModules)
		moduleGroup.GET("/:id", s.dependencies.ModuleController.GetModuleById)
	}
}

func (s *Server) TeamsV1(group *gin.RouterGroup) {
	teamsGroup := group.Group("/teams")
	{
		teamsGroup.Use(s.AdminMiddleware())
		teamsGroup.POST("", s.dependencies.TeamController.CreateTeam)
		teamsGroup.PUT("/:id", s.dependencies.TeamController.UpdateTeam)
		teamsGroup.DELETE("/:id", s.dependencies.TeamController.DeleteTeam)
		teamsGroup.GET("", s.UserMiddleware(), s.dependencies.TeamController.GetTeams)
		teamsGroup.POST("/members", s.dependencies.TeamController.AddTeamMember)
		teamsGroup.DELETE("/members/:id", s.dependencies.TeamController.RemoveTeamMember)
		teamsGroup.GET("/:id/members", s.UserMiddleware(), s.dependencies.TeamController.GetTeamMembers)
	}
}
