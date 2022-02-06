package punching

import (
	"context"

	"github.com/rs/zerolog/log"

	"rendezvous/internal/model"
	"rendezvous/pkg/proto"
)

func (p *Puncher) register(req request) {
	reg := &proto.Registration{}
	if err := reg.Unmarshal(req.data); err != nil {
		log.Error().Err(err).Msg("failure to unmarshal registration")
		return
	}
	user := &model.User{
		Name:          reg.User,
		LocalAddress:  reg.Address,
		GlobalAddress: req.addr.String(),
	}
	if err := p.u.AddUser(context.Background(), user); err != nil {
		log.Error().Err(err).Msg("failure to add new user")
	}
}
