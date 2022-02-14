package client

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Client struct {
	name              string
	rendezvousAddress net.Addr
	localAddress      string
	conn              net.PacketConn
	listener          *Listener
	// we use it to notify upper lever to start serve connections
	signalChan chan []byte
}

const (
	network        = "udp"
	approveTimeout = time.Second * 3
	connectTimeout = time.Second * 5
	punchTimeout   = time.Minute
	clientID       = "client"
	pingMessage    = "ping"
)

func NewClient(ctx context.Context, port string, rendezvousAddress string) (*Client, error) {
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

	cli := &Client{
		rendezvousAddress: addr, localAddress: localAddress, conn: conn, listener: l,
		signalChan: make(chan []byte),
	}

	// starting listen for initiate connection messages from server
	go cli.waitForSignalToStartPunch(ctx)
	return cli, nil
}

func (c *Client) StopListener() {
	c.listener.stop <- struct{}{}
}

func (c *Client) GetConnection() net.PacketConn {
	return c.conn
}

func (c *Client) GetSignalChan() <-chan []byte {
	return c.signalChan
}

func (c *Client) Close() error {
	return c.conn.Close()
}
