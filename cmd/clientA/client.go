package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"rendezvous/pkg/proto/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClient(ctx, "a", "11000", "localhost:9000")
	if err != nil {
		log.Fatal().Err(err).Msg("failure to create client")
	}
	defer cli.Close()
	if err := cli.Register(); err != nil {
		log.Fatal().Err(err).Msgf("failure to register client %s", err)
	}

	<-time.After(time.Minute * 20)
}
