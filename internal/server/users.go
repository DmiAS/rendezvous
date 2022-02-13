package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) getUsers(c *fiber.Ctx) error {
	name := c.Params(userName)
	users := s.u.GetUsers(name)
	return c.Status(http.StatusOK).JSON(users)
}
