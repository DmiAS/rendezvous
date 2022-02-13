package client

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/DmiAS/rendezvous/pkg/proto"
)

// Register send registration request to server
func (c *Client) Register(name string) error {
	c.name = name
	if err := c.sendRegRequest(); err != nil {
		return err
	}

	ch := c.listener.Subscribe(clientID, proto.RegisterApprove)
	defer c.listener.Unsubscribe(clientID, proto.RegisterApprove)
	return c.waitRegApprove(ch)
}

func (c *Client) sendRegRequest() error {
	// prepare data about user
	reg := &proto.Registration{
		User:    c.name,
		Address: c.localAddress,
	}
	header := &proto.Header{Action: proto.RegisterAction}
	packet := &proto.Packet{Header: header, Data: reg}
	data, err := packet.Marshal()
	if err != nil {
		return fmt.Errorf("failure to marshal packet: %s", err)
	}

	// send packet via udp connection
	n, err := c.conn.WriteTo(data, c.rendezvousAddress)
	if err != nil {
		return fmt.Errorf("failure to send data: %s", err)
	}
	log.Debug().Msgf("written %d bytes of %d total", n, len(data))
	return nil
}

// we are waiting for a response from the server, the text of the message is not important to us,
// it is important to receive it in order to confirm registration
func (c *Client) waitRegApprove(ch chan Request) error {
	ctx, cancel := context.WithTimeout(context.Background(), approveTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout")
		case <-ch:
			return nil
		}
	}
}
