package server

import "github.com/gofiber/fiber/v2"

func (s *Server) ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON("Pong...")
}
