package main

import (
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	local, err := net.ResolveUDPAddr("udp", "localhost:9000")
	if err != nil {
		log.Fatal().Err(err).Msg("a")
	}
	// _, err := net.ResolveUDPAddr("udp", "8.8.8.8:9000")
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("b")
	// }
	conn, err := net.ListenUDP("udp", local)
	if err != nil {
		log.Fatal().Err(err).Msg("c")
	}
	b := [512]byte{}
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500))
	if _, _, err := conn.ReadFrom(b[:]); err != nil {
		log.Fatal().Err(err).Msg("d")
	}
}
