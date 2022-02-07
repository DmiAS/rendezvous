package router

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"

	"rendezvous/internal/model"
)

func (r *Router) GetUsers(c *fiber.Ctx) error {
	users, err := r.u.GetUsers(c.UserContext())
	if err != nil {
		log.Error().Err(err).Msg("failure to get users")
		return c.Status(http.StatusInternalServerError).JSON(model.ErrorResponse{Msg: "can not get list of users"})
	}
	return c.Status(http.StatusOK).JSON(users)
}
