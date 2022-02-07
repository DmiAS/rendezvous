package main

import (
	"github.com/rs/zerolog/log"

	"rendezvous/pkg/proto/client"
)

func main() {
	cli, err := client.NewClient("a", "localhost:9000", "localhost:11000")
	if err != nil {
		log.Fatal().Err(err).Msg("failure to create client")
	}
	defer cli.Close()
	if err := cli.Register(); err != nil {
		log.Fatal().Err(err).Msgf("failure to register client %s", err)
	}
}
