package punching

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/internal/model"
)

type UserService interface {
	AddUser(user *model.User)
	BlockUser(user string) error
	GetUser(login string) (*model.User, error)
	DeleteUser(login string)
}

type Puncher struct {
	u        UserService
	requests chan request
	pc       net.PacketConn
}

const (
	network = "udp"
	port    = ":9000"
	workers = 5
)

func NewPuncher(u UserService) *Puncher {
	return &Puncher{u: u, requests: make(chan request, workers)}
}

func (p *Puncher) Listen(ctx context.Context) {
	var err error
	p.pc, err = net.ListenPacket(network, port)
	if err != nil {
		log.Fatal().Err(err).Msg("failure to create socket")
	}
	defer p.pc.Close()

	log.Info().Msgf("server started listen on: %s", p.pc.LocalAddr().String())
	go p.handleActions(ctx)
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("punch server shutdown")
			return
		default:
			p.handleConnection()
		}
	}
}

func (p Puncher) handleConnection() {
	data, clientAddr, err := readAll(p.pc)
	if err != nil {
		log.Error().Err(err).Msg("failure to read datagram from socket")
	}
	req := request{data: data, addr: clientAddr}
	log.Debug().Msgf("sending new request", req)
	p.requests <- req
}

func readAll(pc net.PacketConn) ([]byte, net.Addr, error) {
	buffer := [512]byte{}
	read, addr, err := pc.ReadFrom(buffer[:])
	if err != nil {
		return nil, nil, err
	}
	log.Debug().Msgf("%d bytes was read", read)
	return buffer[:read], addr, nil
}
