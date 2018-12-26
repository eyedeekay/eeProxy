package tunmanager

import (
	"context"
	"fmt"
	"log"
	"net"
)

import (
	"github.com/eyedeekay/eeproxy/conn"
	"github.com/eyedeekay/eeproxy/resolve"
	"github.com/eyedeekay/eeproxy/rewriter"
	"github.com/eyedeekay/go-socks5"
	"github.com/eyedeekay/sam3"
)

type Manager struct {
	resolver.Resolver
	socks5.Config
	*sam3.SAM
	listen  net.Listener
	server  *socks5.Server
	conns   []*conn.Conn
	datadir string
	host    string
	port    string
	samhost string
	samport string
	samopts []string
}

func (m Manager) Serve() error {
	var err error
	if m.listen, err = net.Listen("tcp", m.host+":"+m.port); err != nil {
		return nil, err
	}
	if err := m.server.Serve(m.listen); err != nil {
		return err
	}
	return nil
}

func (m *Manager) DialI2P(ctx context.Context, addr string) (*sam3.SAMConn, error) {
	i2paddr, err := sam3.NewI2PAddrFromString(addr)
	if err != nil {
		return nil, err
	}
	for i, c := range m.conns {
		if i2paddr.Base32() == c.RemoteAddr().(*sam3.I2PAddr).Base32() {
			log.Println("Found destination for address:", i2paddr.Base32(), "at position", i)
			return c.SAMConn, nil
		}
	}
	log.Println("Creating a new connection in connection tree.", m.datadir)
	newconn, err := conn.NewConn(*m.SAM, addr, m.datadir, m.samopts)
	if err != nil {
		return nil, err
	}
	m.conns = append(m.conns, &newconn)
	log.Println("Generated destination for address:", i2paddr.Base32(), "at position", len(m.conns)-1)
	return m.conns[len(m.conns)-1].SAMConn, nil
}

func (m Manager) Dial(ctx context.Context, network, addr string) (net.Conn, error) {
	log.Println("Dialing:", network, addr)
	return m.DialI2P(ctx, addr)
}

func (m Manager) Cleanup() error {
	log.Println("Initiating cleanup.")
	for _, c := range m.conns {
		if err := c.Cleanup(); err != nil {
			return err
		}
	}
	log.Println("Finished cleanup.")
	m.listen.Close()
	return nil
}

func NewManager(host, port, samhost, samport, datadir string, samopts []string) (*Manager, error) {
	return NewManagerFromOptions(
		SetSocksHost(host),
		SetSocksPort(port),
		SetHost(samhost),
		SetPort(samport),
		SetDataDir(datadir),
		SetSAMOpts(samopts),
	)
}

func NewManagerFromOptions(opts ...func(*Manager) error) (*Manager, error) {
	var m Manager
	m.samhost = "127.0.0.1"
	m.samport = "7656"
	m.datadir = "./tunnels/"
	m.host = "127.0.0.1"
	m.port = "7950"
	for _, o := range opts {
		if err := o(&m); err != nil {
			return nil, err
		}
	}
	log.Println("preparing SAM-multiplexing SOCKS proxy on", m.host, m.port, "->", m.samhost, m.samport)
	var err error
	m.SAM, err = sam3.NewSAM(m.samhost + ":" + m.samport)
	if err != nil {
		return nil, err
	}
	if r, err := resolver.NewResolver(m.samhost, m.samport); err == nil {
		m.Config = socks5.Config{
			Resolver: r,
			Dial:     m.Dial,
			Rewriter: rewriter.NewRewriter(),
		}
		return &m, nil
	}
	m.server, err = socks5.New(&m.Config)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("Resolver creation error.")
}
