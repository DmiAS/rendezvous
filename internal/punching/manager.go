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
			header, data, err := proto.GetHeader(request.data)
			if err != nil {
				log.Error().Err(err).Msg("failure to parse request header")
				return
			}
			// change data because we've already read first byte
			request.data = data
			switch header.Action {
			case proto.RegisterAction:
				p.register(request)
			case proto.RequestForConnection:
				p.initConnection(request)
			}
		}
	}
}
