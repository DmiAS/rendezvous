package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"rendezvous/internal/config"
	"rendezvous/internal/router"
)

type Server struct {
	hostAddress string
	app         *fiber.App
	router      *router.Router
}

func NewServer(cfg config.ServerConfig, router *router.Router) *Server {
	srv := &Server{
		hostAddress: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		router:      router,
	}
	srv.app = fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
			CaseSensitive:         true,
			StrictRouting:         true,
		},
	)
	srv.app.Use(recover.New(), logger.New())
	srv.initRoutes()
	return srv
}

func (s *Server) Run() {
	if err := s.app.Listen(s.hostAddress); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msgf("the HTTP rest stopped with unknown error")
	}
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
