package server

const (
	apiPrefix = "/api/v1"
	ping      = "ping"
	users     = "/users"
)

func (s *Server) initRoutes() {
	s.app.Get(apiPrefix+ping, s.router.Ping)
	s.app.Get(apiPrefix+users, s.router.GetUsers)
}
