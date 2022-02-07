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
	name              string
	rendezvousAddress net.Addr
	localAddress      string
	conn              net.PacketConn
	listener          *Listener
}

const (
	network        = "udp"
	approveTimeout = time.Second * 3
	connectTimeout = time.Second * 5
	punchTimeout   = time.Minute
	clientID       = "client"
)

func NewClient(ctx context.Context, name string, rendezvousAddress string, localAddress string) (*Client, error) {
	conn, err := net.ListenPacket(network, localAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to open udp port: %s", err)
	}

	addr, err := net.ResolveUDPAddr(network, rendezvousAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to resolve rendezvous adress %s: %s", rendezvousAddress, err)
	}

	// start listening
	l := NewListener(conn)
	go l.Listen(ctx)

	return &Client{
		rendezvousAddress: addr, localAddress: localAddress, conn: conn, name: name,
		listener: l,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Register() error {
	if err := c.sendRegRequest(); err != nil {
		return err
	}

	return c.waitRegApprove()
}

func (c *Client) sendRegRequest() error {
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

	ch := c.listener.Subscribe(clientID, proto.RegisterApprove)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout")
		case <-ch:
			c.listener.Unsubscribe(clientID, proto.RegisterApprove)
			return nil
		}
	}
}

func (c *Client) ConnectTo(targetName string) (net.Addr, error) {
	request := &proto.ConnRequest{
		Initiator: c.name,
		Target:    targetName,
	}

	data, err := request.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failure to marshal connection request: %s", err)
	}

	n, err := c.conn.WriteTo(data, c.rendezvousAddress)
	if err != nil {
		return nil, fmt.Errorf("failure to send data to server: %s", err)
	}
	log.Debug().Msgf("%d bytes sent from %d", n, len(data))

	connResp, err := c.waitConnResponse()
	if err != nil {
		return nil, fmt.Errorf("failure to get response from server: %s", err)
	}

	addr, err := c.punch(connResp)
	if err != nil {
		return nil, fmt.Errorf("failure to punch: %s", err)
	}
	return addr, nil
}

func (c *Client) waitConnResponse() (*proto.ConnResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	resp := make(chan []byte)
	go c.waitServerResponse(ctx, connectTimeout, resp)
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failure to get registration approve before timeout")
		case data := <-resp:
			connResp := &proto.ConnResponse{}
			if err := connResp.Unmarshal(data); err != nil {
				log.Debug().Err(err).Msg("failure to unmarshal connection response")
				continue
			}
			return connResp, nil
		}
	}
}

func (c *Client) waitServerResponse(ctx context.Context, timeout time.Duration, resp chan []byte) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stop waiting server response")
		default:
			data := [512]byte{}
			if err := c.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
				log.Debug().Err(err).Msg("failure to set timeout")
				return
			}
			n, addr, err := c.conn.ReadFrom(data[:])
			if err != nil {
				log.Error().Err(err).Msg("failure to read from conn")
				return
			}
			if addr.String() == c.rendezvousAddress.String() {
				resp <- data[:n]
				return
			}
		}
	}
}

func (c *Client) punch(connResp *proto.ConnResponse) (net.Addr, error) {
	ping := []byte("ping")

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

	if err := c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 500)); err != nil {
		log.Debug().Err(err).Msg("failure to send timeout")
	}
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failure to punch, timeout")
		default:
			_, err := c.conn.WriteTo(ping, globalAddress)
			if err != nil {
				return globalAddress, nil
			}
			log.Debug().Err(err).Msgf("failed to send data on globalAddress %s", globalAddress)

			_, err = c.conn.WriteTo(ping, localAddress)
			if err != nil {
				return localAddress, nil
			}
			log.Debug().Err(err).Msgf("failed to send data on localAddress %s", localAddress)

			buf := [512]byte{}
			if _, addr, err := c.conn.ReadFrom(buf[:]); err != nil {
				return addr, nil
			}
		}
	}
}
