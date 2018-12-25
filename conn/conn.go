package conn

import (
	//"os"
	"github.com/eyedeekay/sam3"
)

type Conn struct {
	sam3.I2PKeys
	*sam3.StreamSession
	*sam3.SAMConn
	path string
}

func (c Conn) FindKeys() bool {
	return false
}

func (c Conn) SaveKeys() (*sam3.I2PKeys, error) {
	c.I2PKeys, err = sam.NewKeys()
	if err != nil {
		return nil, err
	}
	return &c.I2PKeys, nil
}

func (c Conn) LoadKeys() (*sam3.I2PKeys, error) {
	return &c.I2PKeys, nil
}

func (c Conn) Keys() (*sam3.I2PKeys, error) {
	if c.FindKeys() {
		return c.LoadKeys()
	}
	return c.SaveKeys()
}

func (m Conn) Cleanup() error {
	if err := m.SAMConn.Close(); err != nil {
		return err
	}
	if err := m.StreamSession.Close(); err != nil {
		return err
	}
	return nil
}

func NewConn(sam *sam3.SAM, addr, path string, opts []string) (*Conn, error) {
	var c Conn
	var err error
	c.path = path + addr + ".i2pkeys"
	c.I2PKeys, err = c.Keys()
	if err != nil {
		return nil, err
	}
	c.StreamSession, err = sam.NewStreamSession(c.I2PKeys.Addr().Base32()[0:10], c.I2PKeys, opts)
	if err != nil {
		return nil, err
	}
	i2paddr, err := sam.Lookup(addr)
	if err != nil {
		return nil, err
	}
	c.SAMConn, err = c.StreamSession.DialI2P(i2paddr)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
