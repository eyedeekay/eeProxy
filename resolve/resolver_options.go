package resolver

import (
	"fmt"
	"strconv"
)

//Option is a Resolver option
type Option func(*Resolver) error

//SetHost sets the host of the client's SAM bridge
func SetHost(s string) func(*Resolver) error {
	return func(c *Resolver) error {
		c.samhost = s
		return nil
	}
}

//SetPort sets the port of the client's SAM bridge
func SetPort(v string) func(*Resolver) error {
	return func(c *Resolver) error {
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number.")
		}
		if port < 65536 && port > -1 {
			c.samport = v
			return nil
		}
		return fmt.Errorf("Invalid port.")
	}
}
