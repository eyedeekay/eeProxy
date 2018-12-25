package conn

import (
	//"os"
	"github.com/eyedeekay/sam3"
)

type Conn struct {
	*sam3.SAMConn
	path string
}

func (c Conn) SaveKeys() {

}

func NewConn(conn *sam3.SAMConn, path string) (*Conn, error) {
	return GenConn(conn, path), nil
}

func GenConn(conn *sam3.SAMConn, path string) *Conn {
	return &Conn{SAMConn: conn, path: path}
}
