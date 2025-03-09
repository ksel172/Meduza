package server

import "github.com/gin-contrib/cors"

func (s *Server) RegisterRoutes() {
    s.engine.Use(cors.New(cors.Config{
        AllowOriginFunc: func(origin string) bool {
            return true // Any origin
        },
        AllowCredentials: true,
        AllowMethods:     []string{"GET", "PUT", "PATCH", "DELETE", "POST", "OPTIONS"},
        AllowHeaders:     []string{"Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Credentials"},
    }))

    apiGroup := s.engine.Group("/api")
    apiGroup.Use(s.HandleCors())
    v1Group := apiGroup.Group("/v1")

    s.AuthV1(v1Group)
    s.UsersV1(v1Group)
    s.AgentsV1(v1Group)
    // s.CheckInV1(v1Group)
    s.ListenersV1(v1Group)
    s.PayloadV1(v1Group)
    s.ModuleV1(v1Group)
    s.TeamsV1(v1Group)
}