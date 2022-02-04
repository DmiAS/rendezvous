package router

import "github.com/gofiber/fiber/v2"

func (r *Router) Ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON("Pong...")
}
