package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/models"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"net/http"
)

type DependencyContainer struct {
	UserController      *handlers.UserController
	RedisService        *repos.Service
	AuthController      *handlers.AuthController
	JwtService          *models.JWTService
	AdminController     *handlers.AdminController
	AgentController     *handlers.AgentController
	CheckInController   *handlers.CheckInController
	ListenersController *handlers.ListenerHandler
}

type Server struct {
	host         string
	port         int
	engine       *gin.Engine
	dependencies *DependencyContainer
}

func NewServer(dependencies *DependencyContainer) *Server {

	// Declare Server config
	server := &Server{
		host:         conf.GetMeduzaServerHostname(),
		port:         conf.GetMeduzaServerPort(),
		engine:       gin.Default(),
		dependencies: dependencies,
	}

	server.RegisterRoutes()

	return server
}

func (s *Server) RegisterRoutes() http.Handler {
	router := gin.Default()

	apiGroup := router.Group("/api")
	{
		v1Group := apiGroup.Group("/v1")
		{
			// Authentication Routes
			authRoutes := v1Group.Group("/auth")
			{
				authRoutes.Use(s.HandleCors())
				authRoutes.POST("/register", s.dependencies.UserController.AddUsers)
				authRoutes.POST("/add-admin", s.dependencies.AdminController.CreateAdmin)
				authRoutes.POST("/login", s.dependencies.AuthController.LoginController)
				authRoutes.GET("/refresh-token", s.dependencies.AuthController.RefreshTokenController)
			}

			adminProtectedRoutes := v1Group.Group("/teamserver")
			{
				adminProtectedRoutes.Use(s.HandleCors())
				adminProtectedRoutes.Use(s.AdminMiddleware())
				adminProtectedRoutes.GET("/users", func(context *gin.Context) {
					s.dependencies.UserController.GetUsers(context)
				})
			}

			agentsGroup := v1Group.Group("/agents")
			{
				agentsGroup.GET("/", func(context *gin.Context) {
					s.dependencies.AgentController.GetAgent(context.Writer, context.Request)
				})
				agentsGroup.PUT("/", func(context *gin.Context) {
					s.dependencies.AgentController.UpdateAgent(context.Writer, context.Request)
				})
				agentsGroup.DELETE("/", func(context *gin.Context) {
					s.dependencies.AgentController.DeleteAgent(context.Writer, context.Request)
				})

				// Agent Tasks API
				agentsGroup.GET("/tasks", func(context *gin.Context) {
					s.dependencies.AgentController.GetAgentTasks(context.Writer, context.Request)
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
			}

			checkinGroup := v1Group.Group("/checkin")
			{
				checkinGroup.POST("/", func(context *gin.Context) {
					s.dependencies.CheckInController.CreateAgent(context.Writer, context.Request)
				})
				checkinGroup.GET("/", func(context *gin.Context) {
					s.dependencies.CheckInController.GetTasks(context.Writer, context.Request)
				})
			}

			listenersGroup := v1Group.Group("/listeners")
			{
				listenersGroup.POST("", s.dependencies.ListenersController.CreateListener)
				listenersGroup.GET("/:id", s.dependencies.ListenersController.GetListener)
				listenersGroup.PUT("", s.dependencies.ListenersController.UpdateListener)
				listenersGroup.DELETE(":id", s.dependencies.ListenersController.DeleteListener)
			}
		}
	}

	// Default HelloWorld Handler
	router.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello, World!")
	})

	return router
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.engine.Run(addr)
}
