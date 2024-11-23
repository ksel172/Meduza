package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/conf"
	"github.com/ksel172/Meduza/teamserver/internal/api/handlers"
	"github.com/ksel172/Meduza/teamserver/internal/models"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
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

func (s *Server) RegisterRoutes() {

	apiGroup := s.engine.Group("/api")
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
				adminProtectedRoutes.GET("/users", s.dependencies.UserController.GetUsers)
				adminProtectedRoutes.POST("/users", s.dependencies.UserController.AddUsers)
			}

			agentsGroup := v1Group.Group("/agents")
			{
				agentsGroup.GET("/", s.dependencies.AgentController.GetAgent)
				agentsGroup.PUT("/", s.dependencies.AgentController.UpdateAgent)
				agentsGroup.DELETE("/", s.dependencies.AgentController.DeleteAgent)

				// Agent Tasks API
				agentsGroup.GET("/tasks", s.dependencies.AgentController.GetAgentTasks)
				agentsGroup.POST("/tasks", s.dependencies.AgentController.CreateAgentTask)
				agentsGroup.DELETE("/tasks", s.dependencies.AgentController.DeleteAgentTasks)
				agentsGroup.DELETE("/tasks/task", s.dependencies.AgentController.DeleteAgentTask)
			}

			checkinGroup := v1Group.Group("/checkin")
			{
				checkinGroup.POST("/", s.dependencies.CheckInController.CreateAgent)
				checkinGroup.GET("/", s.dependencies.CheckInController.GetTasks)
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
	s.engine.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello, World!")
	})
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	return s.engine.Run(addr)
}
