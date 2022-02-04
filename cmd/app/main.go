package main

import (
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"rendezvous/internal/config"
	"rendezvous/internal/router"
	"rendezvous/internal/server"
)

const (
	configEnv = "CONFIG_PATH"
)

func main() {
	// read configuration
	cfg, err := config.NewConfig(os.Getenv(configEnv))
	if err != nil {
		panic("failure to read config:" + err.Error())
	}

	// initialize logger
	setupLogger(cfg)

	// initialize repository

	// initialize service

	// initialize router
	r := router.NewRouter()
	// initialize server
	srv := server.NewServer(cfg.Server, r)
	// run server
	runServer(srv)
}

func runServer(srv *server.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Kill, os.Interrupt)
	go srv.Run()
	<-quit
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
	if cfg.Server.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// to enable stack tracing
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
