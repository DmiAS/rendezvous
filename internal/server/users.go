package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) getUsers(c *fiber.Ctx) error {
	users := s.u.GetUsers()
	return c.Status(http.StatusOK).JSON(users)
}
