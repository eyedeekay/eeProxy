package tunmanager

import (
	"context"
	"fmt"
	"log"
	"net"
)

import (
	"github.com/eyedeekay/eeproxy/resolve"
	"github.com/eyedeekay/go-socks5"
	"github.com/eyedeekay/sam3"
)

type Manager struct {
	resolver.Resolver
	socks5.Config
	sam3.StreamSession
	conns []*sam3.SAMConn
}

func (m Manager) Serve() error {
	return nil
}

func (m Manager) DialI2P(ctx context.Context, addr string) (*sam3.SAMConn, error) {
	i2paddr, err := sam3.NewI2PAddrFromString(addr)
	if err != nil {
		return nil, err
	}
	for i, c := range m.conns {
		if i2paddr.Base32() == c.RemoteAddr().(*sam3.I2PAddr).Base32() {
			log.Println("Found destination for address:", i2paddr.Base32(), "at position", i)
			return c, nil
		}
	}
	newconn, err := m.StreamSession.DialI2P(i2paddr)
	if err != nil {
		return nil, err
	}
	m.conns = append(m.conns, newconn)
	log.Println("Generated destination for address:", i2paddr.Base32(), "at position", len(m.conns)-1)
	return m.conns[len(m.conns)-1], nil
}

func (m Manager) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	return m.DialI2P(ctx, addr)
}

func NewManager() (*Manager, error) {
	return NewManagerFromOptions()
}

func NewManagerFromOptions() (*Manager, error) {
	var m Manager
	if r, err := resolver.NewResolver(); err == nil {
		m.Config = socks5.Config{
			Resolver: r,
			Dial:     m.Dial,
		}
		return &m, nil
	}
	return nil, fmt.Errorf("Resolver creation error.")
}
