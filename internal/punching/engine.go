package punching

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/rs/zerolog/log"

	"rendezvous/internal/model"
)

type UserService interface {
	AddUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, login string) error
}

type Puncher struct {
	u        UserService
	requests chan request
}

const (
	network = "udp"
	port    = "9000"
	workers = 5
)

func NewPuncher(u UserService) *Puncher {
	return &Puncher{u: u, requests: make(chan request, workers)}
}

func (p *Puncher) Listen(ctx context.Context) error {
	pc, err := net.ListenPacket(network, port)
	if err != nil {
		return fmt.Errorf("failture to create socket")
	}
	defer pc.Close()

	for {
		p.handleConnection(pc)
	}
}

func (p Puncher) handleConnection(pc net.PacketConn) {
	data, clientAddr, err := readAll(pc)
	if err != nil {
		log.Error().Err(err).Msg("failure to read datagram from socket")
	}
	req := request{data: data, addr: clientAddr}
	log.Debug().Msgf("sending new request: %+v", req)
	p.requests <- req
}

func readAll(pc net.PacketConn) ([]byte, net.Addr, error) {
	buffer := [512]byte{}
	res := &bytes.Buffer{}
	var clientAddr net.Addr
	for {
		read, addr, err := pc.ReadFrom(buffer[:])
		if err != nil {
			if err == io.EOF {
				clientAddr = addr
				break
			}
			return nil, nil, err
		}
		res.Write(buffer[:read])
	}
	return res.Bytes(), clientAddr, nil
}
