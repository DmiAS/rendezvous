package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/DmiAS/rendezvous/internal/config"
	"github.com/DmiAS/rendezvous/internal/punching"
	"github.com/DmiAS/rendezvous/internal/repository/mem"
	"github.com/DmiAS/rendezvous/internal/server"
	"github.com/DmiAS/rendezvous/internal/service"
)

func main() {
	// read configuration
	cfg := config.NewConfig()

	// initialize logger
	setupLogger(cfg)

	// initialize repository
	repo := mem.NewUsersRepository()

	// initialize service
	userService := service.NewUserService(repo)

	// initialize server
	srv := server.NewServer(cfg, userService)

	// initialize punch server
	punchSrv := punching.NewPuncher(userService)

	// run server
	runServer(srv, punchSrv)
}

func runServer(srv *server.Server, punchServer *punching.Puncher) {
	// create context
	ctx, cancel := context.WithCancel(context.Background())

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Kill, os.Interrupt)
	go srv.Run()
	go punchServer.Listen(ctx)
	<-quit
	cancel()

	log.Info().Msg("server shutting down")
	if err := srv.Shutdown(); err != nil {
		log.Fatal().
			Err(err).
			Msg("failure to shutdown server")
	}
	log.Info().
		Msg("server stopped")
}

func setupLogger(cfg *config.Config) {
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// to enable stack tracing
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
