package client

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/DmiAS/rendezvous/pkg/proto"
)

type Request struct {
	addr net.Addr
	data []byte
}

type subscriber struct {
	name string
	ch   chan Request
}

type Subscribers map[uint8][]subscriber

type Listener struct {
	mu   sync.RWMutex
	subs Subscribers
	conn net.PacketConn
}

const (
	sendTimeout = time.Second
	chanSize    = 5
)

func NewListener(conn net.PacketConn) *Listener {
	return &Listener{conn: conn, mu: sync.RWMutex{}, subs: make(Subscribers)}
}

func (l *Listener) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("listener stopped")
		default:
			buffer := [512]byte{}
			n, addr, err := l.conn.ReadFrom(buffer[:])
			if err != nil {
				log.Error().Err(err).Msg("failure to read data from conn")
			}
			l.notifySubs(
				Request{
					addr: addr,
					data: buffer[:n],
				},
			)
		}
	}
}

func (l *Listener) Subscribe(name string, event uint8) chan Request {
	ch := make(chan Request, chanSize)
	sub := subscriber{name: name, ch: ch}
	l.mu.Lock()
	if subs, ok := l.subs[event]; !ok {
		l.subs[event] = []subscriber{sub}
	} else {
		l.subs[event] = append(subs, sub)
	}
	l.mu.Unlock()
	return ch
}

func (l *Listener) Unsubscribe(name string, event uint8) {
	l.mu.Lock()
	if subs, ok := l.subs[event]; ok {
		for i := range subs {
			if subs[i].name == name {
				subs[i] = subs[len(subs)-1]
				log.Debug().Msgf("subs before unsub = %+v", subs)
				subs = subs[:len(subs)-1]
				log.Debug().Msgf("subs after unsub = %+v", subs)
				break
			}
		}
		l.subs[event] = subs
	}
	l.mu.Unlock()
}

func (l *Listener) notifySubs(req Request) {
	header, data, err := proto.GetHeader(req.data)
	if err != nil {
		log.Error().Err(err).Msg("failure to extract header")
		return
	}
	req.data = data
	log.Debug().Msgf("new message with action: %d", header)

	l.mu.RLock()
	for _, sub := range l.subs[header.Action] {
		ctx, cancel := context.WithTimeout(context.Background(), sendTimeout)
		select {
		case <-ctx.Done():
			log.Debug().Msgf("failure to notify subscriber for event %d", header.Action)
		case sub.ch <- req:
			log.Debug().Msgf("notify sub with event %d", header.Action)
		}
		cancel()
	}
	l.mu.RUnlock()
}
