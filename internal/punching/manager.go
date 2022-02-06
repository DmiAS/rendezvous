package punching

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"rendezvous/pkg/proto"
)

type request struct {
	data []byte
	addr net.Addr
}

func (p *Puncher) handleActions(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stop handling punching requests")
			return
		case request := <-p.requests:
			header := &proto.Header{}
			if err := header.Unmarshal(request.data); err != nil {
				log.Error().Err(err).Msg("failure to parse request header")
			}
			// change data because we've already read first byte
			request.data = request.data[1:]
			switch header.Action {
			case proto.RegisterAction:
				p.register(request)
			case proto.InitiateConnectionAction:
				p.initConnection(request)
			}
		}
	}
}
