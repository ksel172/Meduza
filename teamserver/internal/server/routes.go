package server

import "github.com/gin-contrib/cors"

func (s *Server) RegisterRoutes() {
	// Routes will look like this :- /api/v1
	s.engine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "PUT", "PATCH", "DELETE", "POST"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
	}))

	apiGroup := s.engine.Group("/api")
	/* 	apiGroup.Use(s.HandleCors()) */
	v1Group := apiGroup.Group("/v1")

	s.AuthV1(v1Group)
	s.AdminV1(v1Group)
	s.AgentsV1(v1Group)
	// s.CheckInV1(v1Group)
	s.ListenersV1(v1Group)
	s.PayloadV1(v1Group)
	s.ModuleV1(v1Group)
}
