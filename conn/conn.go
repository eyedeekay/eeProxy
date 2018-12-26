package conn

import (
	"bufio"
	"log"
	"math/rand"
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
			os.MkdirAll(c.path, 0755)
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
	log.Println("Cleaning up client connection.", m.name)
	if err := m.SAMConn.Close(); err != nil {
		return err
	}
	log.Println("Cleaning up client session.", m.name)
	if err := m.StreamSession.Close(); err != nil {
		return err
	}
	return nil
}

func NewConn(sam sam3.SAM, addr, path string, opts []string) (*Conn, error) {
	var c Conn
	var err error
	c.SAM = &sam
	c.path = path
	t32, err := sam3.NewI2PAddrFromString(addr)
	c.name = t32.Base32() + ".i2pkeys"
	c.I2PKeys, err = c.Keys()
	if err != nil {
		return nil, err
	}
	c.StreamSession, err = c.SAM.NewStreamSession(c.I2PKeys.Addr().Base32()[0:10]+"-"+RandTunName(), c.I2PKeys, opts)
	if err != nil {
		return nil, err
	}
	i2paddr, err := c.SAM.Lookup(addr)
	if err != nil {
		return nil, err
	}
	c.SAMConn, err = c.StreamSession.DialI2P(i2paddr)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// RandTunName generates a random tunnel names to avoid collisions
func RandTunName() string {
	b := make([]byte, 4)
	for i := range b {
		b[i] = "abcdefghijklmnopqrstuvwxyz"[rand.Intn(len("abcdefghijklmnopqrstuvwxyz"))]
	}
	return string(b)
}
