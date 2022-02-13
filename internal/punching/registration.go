package punching

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/internal/model"
	"github.com/DmiAS/rendezvous/pkg/proto"
)

func (p *Puncher) register(req request) {
	reg := &proto.Registration{}
	if err := reg.Unmarshal(req.data); err != nil {
		msg := "failure to unmarshal registration"
		log.Error().Err(err).Msg(msg)
		p.sendRegisterApprove(req.addr, msg)
		return
	}
	user := &model.User{
		Name:          reg.User,
		LocalAddress:  reg.Address,
		GlobalAddress: req.addr.String(),
	}
	if err := p.u.AddUser(context.Background(), user); err != nil {
		msg := "failure to add new user"
		log.Error().Err(err).Msg(msg)
		p.sendRegisterApprove(req.addr, msg)
		return
	}
	p.sendRegisterApprove(req.addr, "")
}

func (p *Puncher) sendRegisterApprove(addr net.Addr, msg string) {
	approve := &proto.RegistrationApprove{
		Error: false,
		Msg:   msg,
	}
	if msg != "" {
		approve.Error = true
	}
	header := &proto.Header{Action: proto.RegisterApprove}
	packet := &proto.Packet{Data: approve, Header: header}

	data, err := packet.Marshal()
	if err != nil {
		log.Error().Err(err).Msgf("failure to marshal approve information for %s", addr)
	}
	if err := p.send(addr, data); err != nil {
		log.Error().Err(err).Msgf("failure to send approve information for %s", addr)
	}
}
