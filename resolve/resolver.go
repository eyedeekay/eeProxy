package resolver

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

import (
	"github.com/eyedeekay/sam3"
)

type Resolver struct {
	*sam3.SAMResolver
	samhost         string
	samport         string
	allowedSuffixes []string
}

func (r Resolver) Resolve(ctx context.Context, name string) (context.Context, net.Addr, error) {
	return r.ResolveI2P(ctx, name)
}

func (r Resolver) ResolveI2P(ctx context.Context, name string) (context.Context, *sam3.I2PAddr, error) {
	if !r.ValidateI2PAddr(name) {
		return ctx, nil, fmt.Errorf("Error, not an allowed suffix")
	}
	raddr, err := r.SAMResolver.Resolve(name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, &raddr, nil
}

func (r Resolver) Cleanup() error {
	return r.SAMResolver.SAM.Close()
}

func (r Resolver) ValidateI2PAddr(name string) bool {
	noi2p := false
	for _, suffix := range r.allowedSuffixes {
		if strings.HasSuffix(name, suffix) {
			if suffix == ".b32.i2p" {
				if len(name) != 60 {
					noi2p = true
					break
				}
			}
			noi2p = true
		}
	}
	return noi2p
}

func NewResolver(samhost, samport string) (*Resolver, error) {
	return NewResolverFromOptions(SetHost(samhost), SetPort(samport))
}

func NewResolverFromOptions(opts ...func(*Resolver) error) (*Resolver, error) {
	var r Resolver
	r.allowedSuffixes = []string{".i2p", ".b32.i2p"}
	r.samhost = "127.0.0.1"
	r.samport = "7656"
	for _, o := range opts {
		if err := o(&r); err != nil {
			return nil, err
		}
	}
	log.Println("Preparing i2p address resolver")
	var err error
	r.SAMResolver, err = sam3.NewFullSAMResolver("127.0.0.1:7656")
	if err != nil {
		return nil, err
	}
	return &r, nil
}
