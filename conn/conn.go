package conn

import (
	"bufio"
	"os"
	"path/filepath"
)

import (
	"github.com/eyedeekay/sam3"
)

type Conn struct {
	*sam3.SAM
	sam3.I2PKeys
	*sam3.StreamSession
	*sam3.SAMConn
	path string
	name string
}

func (c Conn) FindKeys() bool {
	if _, err := os.Stat(c.Path()); os.IsNotExist(err) {
		if _, err := os.Stat(c.path); os.IsNotExist(err) {
			os.MkdirAll(c.path, os.ModeDir)
		}
		return false
	}
	return true
}

func (c Conn) Path() string {
	p := filepath.Join(c.path, c.name)
	return p
}

func (c Conn) SaveKeys() (sam3.I2PKeys, error) {
	var err error
	c.I2PKeys, err = c.SAM.NewKeys()
	if err != nil {
		return sam3.I2PKeys{}, err
	}
	f, err := os.Create(c.Path())
	if err != nil {
		return sam3.I2PKeys{}, err
	}
	defer f.Close()
	filewriter := bufio.NewWriter(f)
	err = sam3.StoreKeysIncompat(c.I2PKeys, filewriter)
	if err != nil {
		return sam3.I2PKeys{}, err
	}
	filewriter.Flush()
	return c.I2PKeys, nil
}

func (c Conn) LoadKeys() (sam3.I2PKeys, error) {
	var err error
	f, err := os.Open(c.Path())
	if err != nil {
		return sam3.I2PKeys{}, err
	}
	defer f.Close()
	filereader := bufio.NewReader(f)
	c.I2PKeys, err = sam3.LoadKeysIncompat(filereader)
	if err != nil {
		return sam3.I2PKeys{}, err
	}
	return c.I2PKeys, nil
}

func (c Conn) Keys() (sam3.I2PKeys, error) {
	if c.FindKeys() {
		return c.LoadKeys()
	}
	return c.SaveKeys()
}

func (m Conn) Cleanup() error {
	if err := m.SAMConn.Close(); err != nil {
		return err
	}
	//if err := m.StreamSession.Close(); err != nil {
	//return err
	//}
	return nil
}

func NewConn(sam *sam3.SAM, addr, path string, opts []string) (*Conn, error) {
	var c Conn
	var err error
	c.SAM = sam
	c.path = path
	c.name = addr + ".i2pkeys"
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
