package punching

import (
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/internal/model"
	"github.com/DmiAS/rendezvous/pkg/proto"
)

const (
	writeDeadline = 5 * time.Second
)

func (p *Puncher) initConnection(req request) {
	connRequest := &proto.ConnRequest{}
	if err := connRequest.Unmarshal(req.data); err != nil {
		log.Error().Err(err).Msg("failure to parse connRequest")
		return
	}

	initiatorData, err := p.u.GetUser(connRequest.Initiator)
	if err != nil {
		log.Error().Err(err).Msgf("failure to get user %s data", connRequest.Initiator)
		return
	}

	targetData, err := p.u.GetUser(connRequest.Target)
	if err != nil {
		log.Error().Err(err).Msgf("failure to get user %s data", connRequest.Target)
		return
	}

	if err := p.sendUserData(initiatorData.GlobalAddress, targetData); err != nil {
		log.Error().
			Err(err).
			Msgf(
				"failure to send to user %s data about %s",
				initiatorData.Name,
				targetData.Name,
			)
		return
	}

	if err := p.sendUserData(targetData.GlobalAddress, initiatorData); err != nil {
		log.Error().
			Err(err).
			Msgf(
				"failure to send to user %s data about %s",
				targetData.Name,
				initiatorData.Name,
			)
		return
	}

	// set block status to each user
	if err := p.u.BlockUser(initiatorData.Name); err != nil {
		log.Error().Err(err).Msgf("failure to block user: %s", initiatorData.Name)
	}

	if err := p.u.BlockUser(targetData.Name); err != nil {
		log.Error().Err(err).Msgf("failure to block user: %s", targetData.Name)
	}
}

func (p *Puncher) sendUserData(globalAddr string, user *model.User) error {
	info := &proto.ConnResponse{GlobalAddress: user.GlobalAddress, LocalAddress: user.LocalAddress}
	header := &proto.Header{Action: proto.ResponseConnection}

	data, err := proto.Packet{Header: header, Data: info}.Marshal()
	if err != nil {
		return fmt.Errorf("failure to marshal conn response data: %+v: %s", info, err)
	}

	addr, err := net.ResolveUDPAddr(network, globalAddr)
	if err != nil {
		return fmt.Errorf("failure to resolve addr: %s : %s", addr, err)
	}

	if err := p.send(addr, data); err != nil {
		return fmt.Errorf("failure to send data to addr: %s: %s", globalAddr, err)
	}
	return nil
}

func (p *Puncher) send(addr net.Addr, data []byte) error {
	log.Debug().Msgf("send data to addr: %s", addr.String())
	n, err := p.pc.WriteTo(data, addr)
	if err != nil {
		return fmt.Errorf("failure to write data to socket: %s", err)
	}
	log.Debug().Msgf("data's len = %d, written = %d", len(data), n)
	return nil
}
