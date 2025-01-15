package server

func (s *Server) RegisterRoutes() {
	// Routes will look like this :- /api/v1
	apiGroup := s.engine.Group("/api")
	apiGroup.Use(s.HandleCors())
	v1Group := apiGroup.Group("/v1")

	s.AuthV1(v1Group)
	s.AdminV1(v1Group)
	s.AgentsV1(v1Group)
	// s.CheckInV1(v1Group)
	s.ListenersV1(v1Group)
	s.PayloadV1(v1Group)
}
