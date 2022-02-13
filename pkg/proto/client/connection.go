package client

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/pkg/proto"
)

// ConnectTo send request for connection to target with particular name
func (c *Client) ConnectTo(targetName string) (net.Addr, error) {
	// sending request to server
	if err := c.initConnection(targetName); err != nil {
		return nil, err
	}

	// waiting for response from server with both addresses of target
	ch := c.listener.Subscribe(clientID, proto.ResponseForInitiator)
	connResp, err := c.waitConnResponse(ch)
	if err != nil {
		c.listener.Unsubscribe(clientID, proto.ResponseForInitiator)
		return nil, fmt.Errorf("failure to get response from server: %s", err)
	}
	c.listener.Unsubscribe(clientID, proto.ResponseForInitiator)

	// initiate punching
	addr, err := c.punch(connResp, proto.PunchInitiatorMessage)
	if err != nil {
		return nil, fmt.Errorf("failure to punch: %s", err)
	}
	return addr, nil
}

func (c *Client) initConnection(targetName string) error {
	request := &proto.ConnRequest{
		Initiator: c.name,
		Target:    targetName,
	}
	header := &proto.Header{Action: proto.RequestForConnection}

	packet := &proto.Packet{Header: header, Data: request}
	data, err := packet.Marshal()
	if err != nil {
		return fmt.Errorf("failure to marshal connection request: %s", err)
	}

	n, err := c.conn.WriteTo(data, c.rendezvousAddress)
	if err != nil {
		return fmt.Errorf("failure to send data to server: %s", err)
	}
	log.Debug().Msgf("%d bytes sent from %d", n, len(data))
	return nil
}

// we are waiting for message from the server to get both addresses of the target
func (c *Client) waitConnResponse(ch chan Request) (*proto.ConnResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failure to get connection before timeout")
		case req := <-ch:
			connResp := &proto.ConnResponse{}
			if err := connResp.Unmarshal(req.data); err != nil {
				log.Debug().Err(err).Msg("failure to unmarshal connection response")
				continue
			}
			return connResp, nil
		}
	}
}
