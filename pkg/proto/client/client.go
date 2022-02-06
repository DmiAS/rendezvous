package client

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/zerolog/log"

	"rendezvous/pkg/proto"
)

type Client struct {
	rendezvousAddress net.Addr
	localAddress      string
	conn              net.PacketConn
}

const (
	network        = "udp"
	approveTimeout = time.Second * 3
)

func NewClient(rendezvousAddress string, localAddress string) (*Client, error) {
	conn, err := net.ListenPacket(network, localAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to open udp port: %s", err)
	}

	addr, err := net.ResolveUDPAddr(network, rendezvousAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to resolve rendezvous adress %s: %s", rendezvousAddress, err)
	}
	return &Client{rendezvousAddress: addr, localAddress: localAddress, conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Register(name string) error {
	if err := c.sendRegRequest(name); err != nil {
		return err
	}

	return c.waitRegApprove()
}

func (c *Client) sendRegRequest(name string) error {
	reg := &proto.Registration{
		User:    name,
		Address: c.localAddress,
	}
	data, err := reg.Marshal()
	if err != nil {
		return fmt.Errorf("failure to marshal reg info: %s", err)
	}
	n, err := c.conn.WriteTo(data, c.rendezvousAddress)
	if err != nil {
		return fmt.Errorf("failure to send data: %s", err)
	}
	log.Debug().Msgf("written %d bytes of %d total", n, len(data))
	return nil
}

func (c *Client) waitRegApprove() error {
	ctx, cancel := context.WithTimeout(context.Background(), approveTimeout)
	defer cancel()

	resp := make(chan []byte)
	go func() {
		data := [512]byte{}
		n, addr, err := c.conn.ReadFrom(data[:])
		if err != nil {
			log.Error().Err(err).Msg("failure to read from conn")
		}
		if addr.String() == c.rendezvousAddress.String() {
			resp <- data[:n]
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("failure to get registration approve before timeout")
		case data := <-resp:
			approve := &proto.RegistrationApprove{}
			if err := approve.Unmarshal(data); err != nil {
				log.Debug().Err(err).Msg("failure to unmarshal registration approval")
				continue
			}
			if approve.Error {
				return fmt.Errorf("failure to register: %s", approve.Error)
			}
			return nil
		}
	}
}

// func (c *Client) ConnectTo(name string) (chan []byte, error) {
//
// }
