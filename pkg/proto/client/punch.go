package client

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/pkg/proto"
)

// starts the process of breaking through,
// which consists in sequentially sending messages to each of the addresses until a response is received from at
// least one
func (c *Client) punch(connResp *proto.ConnResponse, action uint8) (net.Addr, error) {
	ping := &proto.Punch{Msg: pingMessage}
	header := &proto.Header{Action: action}

	data, err := proto.Packet{Header: header, Data: ping}.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failure to marshal packet")
	}

	globalAddress, err := net.ResolveUDPAddr(network, connResp.GlobalAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to resolve global address: %s", err)
	}

	localAddress, err := net.ResolveUDPAddr(network, connResp.LocalAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to resolve local address: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), punchTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failure to punch, timeout")
		default:
			_, err := c.conn.WriteTo(data, globalAddress)
			if err == nil {
				return globalAddress, nil
			}
			log.Debug().Err(err).Msgf("failed to send data on globalAddress %s", globalAddress)

			_, err = c.conn.WriteTo(data, localAddress)
			if err == nil {
				return localAddress, nil
			}
			log.Debug().Err(err).Msgf("failed to send data on localAddress %s", localAddress)
		}
	}
}

// here we are waiting for a message from the server that someone wants to join us,
// this is necessary to start the process of breaking through from the target side
func (c *Client) waitForSignalToStartPunch(ctx context.Context) {
	ch := c.listener.Subscribe(c.name, proto.ResponseForTarget)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("stop listen for punching signals")
			return
		case req := <-ch:
			// connResp included initiator's addresses
			connResp := &proto.ConnResponse{}
			if err := connResp.Unmarshal(req.data); err != nil {
				log.Debug().Err(err).Msg("failure to unmarshal connection response")
				continue
			}
			if _, err := c.punch(connResp, proto.PunchTargetMessage); err != nil {
				log.Debug().Err(err).Msg("failure to punch initiator")
			} else {
				// if the breakout was successful, then we send a signal that you need to start listening
				c.signalChan <- struct{}{}
			}
		}
	}
}
