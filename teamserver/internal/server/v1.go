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
		authRoutes.GET("/refresh", s.dependencies.AuthController.RefreshTokenController)
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
		// Agents middleware
		agentsGroup.Use(s.UserMiddleware())

		// Base agent operations
		agentsGroup.GET(fmt.Sprintf("/:%s", models.ParamAgentID), s.dependencies.AgentController.GetAgent)
		agentsGroup.PUT(fmt.Sprintf("/:%s", models.ParamAgentID), s.dependencies.AgentController.UpdateAgent)
		agentsGroup.DELETE(fmt.Sprintf("/:%s", models.ParamAgentID), s.dependencies.AgentController.DeleteAgent)

		// Agent Tasks API
		agentsGroup.GET("/tasks", s.dependencies.AgentController.GetAgentTasks)
		agentsGroup.POST(fmt.Sprintf("/tasks/:%s", models.ParamAgentID), s.dependencies.AgentController.CreateAgentTask)
		agentsGroup.DELETE(fmt.Sprintf("/:%s/tasks", models.ParamAgentID), s.dependencies.AgentController.DeleteAgentTasks)
		agentsGroup.DELETE(fmt.Sprintf("/:%s/tasks/:%s", models.ParamAgentID, models.ParamTaskID), s.dependencies.AgentController.DeleteAgentTask)

		// Agent config API
		agentsGroup.POST(fmt.Sprintf("/config/:%s", models.ParamAgentID), s.dependencies.AgentController.CreateAgentConfig)
		agentsGroup.PUT(fmt.Sprintf("/config/:%s", models.ParamAgentID), s.dependencies.AgentController.UpdateAgentConfig)
		agentsGroup.GET(fmt.Sprintf("/config/:%s", models.ParamAgentID), s.dependencies.AgentController.GetAgentConfig)
		agentsGroup.DELETE(fmt.Sprintf("/config/:%s", models.ParamAgentID), s.dependencies.AgentController.DeleteAgentConfig)

		// Agent info API
		agentsGroup.POST("/info", s.dependencies.AgentController.CreateAgentInfo)
		agentsGroup.PUT(fmt.Sprintf("/info/:%s", models.ParamAgentID), s.dependencies.AgentController.UpdateAgentInfo)
		agentsGroup.GET(fmt.Sprintf("/info/:%s", models.ParamAgentID), s.dependencies.AgentController.GetAgentInfo)
		agentsGroup.DELETE(fmt.Sprintf("/info/:%s", models.ParamAgentID), s.dependencies.AgentController.DeleteAgentInfo)
	}
}
func (s *Server) ListenersV1(group *gin.RouterGroup) {

	listenersGroup := group.Group("/listeners")
	{
		// Listeners middleware
		listenersGroup.Use(s.UserMiddleware())

		// Listener CRUD operations and status info
		listenersGroup.POST("", s.dependencies.ListenerController.CreateListener)
		listenersGroup.GET(fmt.Sprintf("/:%s", models.ParamListenerID), s.dependencies.ListenerController.GetListenerById)
		listenersGroup.GET("/all", s.dependencies.ListenerController.GetAllListeners)
		listenersGroup.PUT(fmt.Sprintf("/:%s", models.ParamListenerID), s.dependencies.ListenerController.UpdateListener)
		listenersGroup.DELETE(fmt.Sprintf("/:%s", models.ParamListenerID), s.dependencies.ListenerController.DeleteListener)
		listenersGroup.GET(fmt.Sprintf("/:%s/status", models.ParamListenerID), s.dependencies.ListenerController.CheckRunningListener)

		// Listener operations
		listenersGroup.POST(fmt.Sprintf("/:%s/start", models.ParamListenerID), s.dependencies.ListenerController.StartListener)
		listenersGroup.POST(fmt.Sprintf("/:%s/stop", models.ParamListenerID), s.dependencies.ListenerController.StopListener)
	}
}
func (s *Server) PayloadV1(group *gin.RouterGroup) {

	payloadsGroup := group.Group("/payloads")
	{
		// Payload CRUD operations and download
		payloadsGroup.Use(s.UserMiddleware())
		payloadsGroup.POST("/", s.dependencies.PayloadController.CreatePayload)
		payloadsGroup.GET("/all", s.dependencies.PayloadController.GetAllPayloads)
		payloadsGroup.DELETE(fmt.Sprintf("/delete/:%s", models.ParamPayloadID), s.dependencies.PayloadController.DeletePayload)
		payloadsGroup.DELETE("/delete/all", s.dependencies.PayloadController.DeleteAllPayloads)
		payloadsGroup.GET(fmt.Sprintf("/download/:%s", models.ParamPayloadID), s.dependencies.PayloadController.DownloadPayload)
	}
}

func (s *Server) ModuleV1(group *gin.RouterGroup) {

	moduleGroup := group.Group("/modules")
	{
		moduleGroup.Use(s.UserMiddleware())
		moduleGroup.POST("/upload", s.dependencies.ModuleController.UploadModule)
		moduleGroup.DELETE(fmt.Sprintf("/delete/:%s", models.ParamModuleID), s.dependencies.ModuleController.DeleteModule)
		moduleGroup.DELETE("/delete/all", s.dependencies.ModuleController.DeleteAllModules)
		moduleGroup.GET("/all", s.dependencies.ModuleController.GetAllModules)
		moduleGroup.GET(fmt.Sprintf("/:%s", models.ParamPayloadID), s.dependencies.ModuleController.GetModuleById)
	}
}

func (s *Server) TeamsV1(group *gin.RouterGroup) {
	teamsGroup := group.Group("/teams")
	{
		teamsGroup.Use(s.AdminMiddleware())
		teamsGroup.POST("", s.dependencies.TeamController.CreateTeam)
		teamsGroup.PUT(fmt.Sprintf("/:%s", models.ParamTeamID), s.dependencies.TeamController.UpdateTeam)
		teamsGroup.DELETE(fmt.Sprintf("/:%s", models.ParamTeamID), s.dependencies.TeamController.DeleteTeam)
		teamsGroup.GET("", s.UserMiddleware(), s.dependencies.TeamController.GetTeams)

		// Team members endpoints
		teamsGroup.POST("/members", s.dependencies.TeamController.AddTeamMember)
		teamsGroup.DELETE(fmt.Sprintf("/members/:%s", models.ParamMemberID), s.dependencies.TeamController.RemoveTeamMember)
		teamsGroup.GET(fmt.Sprintf("/:%s/members", models.ParamTeamID), s.UserMiddleware(), s.dependencies.TeamController.GetTeamMembers)
	}
}
