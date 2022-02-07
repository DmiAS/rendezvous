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
	pingMessage    = "ping"
)

func NewClient(ctx context.Context, name string, port string, rendezvousAddress string) (*Client, error) {
	localAddress, err := resolveLocalIpAddress()
	if err != nil {
		return nil, fmt.Errorf("failure to get local address: %s", err)
	}
	localAddress += ":" + port
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

func resolveLocalIpAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("pc is not connected to the network")
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Register() error {
	if err := c.sendRegRequest(); err != nil {
		return err
	}

	ch := c.listener.Subscribe(clientID, proto.RegisterApprove)
	defer c.listener.Unsubscribe(clientID, proto.RegisterApprove)
	return c.waitRegApprove(ch)
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

func (c *Client) ConnectTo(targetName string) (net.Addr, error) {
	if err := c.initConnection(targetName); err != nil {
		return nil, err
	}

	ch := c.listener.Subscribe(clientID, proto.ResponseConnection)
	connResp, err := c.waitConnResponse(ch)
	if err != nil {
		c.listener.Unsubscribe(clientID, proto.ResponseConnection)
		return nil, fmt.Errorf("failure to get response from server: %s", err)
	}
	c.listener.Unsubscribe(clientID, proto.ResponseConnection)

	addr, err := c.punch(connResp)
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
func (c *Client) punch(connResp *proto.ConnResponse) (net.Addr, error) {
	ping := &proto.Punch{Msg: pingMessage}
	header := &proto.Header{Action: proto.PunchMessage}

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
