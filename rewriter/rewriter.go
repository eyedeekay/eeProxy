package rewriter

import (
	"context"
)
import (
	"github.com/eyedeekay/go-socks5"
	"github.com/eyedeekay/sam3"
)

type Rewriter struct {
	network string
}

func (r Rewriter) Rewrite(ctx context.Context, request *socks5.Request) (context.Context, *socks5.AddrSpec) {
	addr := request.DestAddr
	addr.FQDN = request.DestAddr.ADDR.(*sam3.I2PAddr).Base32()
	log.Println("Correcting FQDN to base32 address.", addr.FQDN)
	return ctx, addr
}

func NewRewriter() *Rewriter {
	var r Rewriter
	return &r
}
