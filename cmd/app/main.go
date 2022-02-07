package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"rendezvous/internal/config"
	"rendezvous/internal/punching"
	"rendezvous/internal/repository/mdb"
	"rendezvous/internal/router"
	"rendezvous/internal/server"
	"rendezvous/internal/service"
	"rendezvous/pkg/mongodb"
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

	// initialize db connection
	mongoConn, err := mongodb.NewConnection(cfg.Db.DSN, cfg.Db.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failure to initialize connection to mongo")
	}
	defer mongoConn.Conn.Disconnect(context.Background())

	// initialize repository
	repo := mdb.NewUsersRepository(mongoConn.Db)

	// initialize service
	userService := service.NewUserService(repo)

	// initialize router
	r := router.NewRouter(userService)

	// initialize server
	srv := server.NewServer(cfg.Server, r)

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
	if cfg.Server.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	// to enable stack tracing
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
