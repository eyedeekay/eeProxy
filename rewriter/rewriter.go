package rewriter

import (
	"context"
	"log"
)
import (
	"github.com/eyedeekay/go-socks5"
	"github.com/eyedeekay/sam3"
)

type Rewriter struct {
	network string
}

func (r Rewriter) Rewrite(ctx context.Context, request *socks5.Request) (context.Context, *socks5.AddrSpec) {
	var addr *socks5.AddrSpec
	switch request.DestAddr.ADDR.(type) {
	case *sam3.I2PAddr:
		addr = request.DestAddr
		//addr.FQDN = request.DestAddr.ADDR.(*sam3.I2PAddr).Base32()
		log.Println("Checking FQDN", addr.FQDN)
	default:
		log.Println(request.DestAddr.String())
		return ctx, &socks5.AddrSpec{}
	}
	return ctx, addr
}

func NewRewriter() *Rewriter {
	var r Rewriter
	return &r
}
