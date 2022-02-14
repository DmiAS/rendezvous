package server

import "github.com/gofiber/fiber/v2/middleware/cors"

const (
	apiPrefix = "/api/v1"
	ping      = "/ping"
	users     = "/users"
	userName  = "user"
)

func (s *Server) initRoutes() {
	api := s.app.Group(apiPrefix).Use(cors.New())
	{
		api.Get(ping, s.ping)
		api.Get(users+"/:"+userName, s.getUsers)
	}
}
