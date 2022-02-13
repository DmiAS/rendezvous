package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/internal/config"
	"github.com/DmiAS/rendezvous/internal/router"
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
			DisableStartupMessage: false,
			CaseSensitive:         true,
			StrictRouting:         true,
		},
	)
	srv.app.Use(recover.New(), logger.New())
	srv.initRoutes()
	return srv
}

func (s *Server) Run() {
	log.Info().Msgf("server started listen on address: %s", s.hostAddress)
	if err := s.app.Listen(s.hostAddress); err != nil && err != http.ErrServerClosed {
		log.Fatal().
			Err(err).
			Msgf("the HTTP rest stopped with unknown error")
	}
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
