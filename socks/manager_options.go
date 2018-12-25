package tunmanager

import (
	"fmt"
	"strconv"
)

//Option is a Manager option
type Option func(*Manager) error

//SetHost sets the host of the client's SAM bridge
func SetHost(s string) func(*Manager) error {
	return func(c *Manager) error {
		c.samhost = s
		return nil
	}
}

//SetDataDir sets the directory to save per-site keys
func SetDataDir(s string) func(*Manager) error {
	return func(c *Manager) error {
		c.datadir = s
		return nil
	}
}

//SetPort sets the port of the client's SAM bridge
func SetPort(v string) func(*Manager) error {
	return func(c *Manager) error {
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

//SetSocksHost sets the host of the client's SAM bridge
func SetHost(s string) func(*Manager) error {
	return func(c *Manager) error {
		c.host = s
		return nil
	}
}

//SetSocksPort sets the port of the client's SAM bridge
func SetPort(v string) func(*Manager) error {
	return func(c *Manager) error {
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number.")
		}
		if port < 65536 && port > -1 {
			c.port = v
			return nil
		}
		return fmt.Errorf("Invalid port.")
	}
}

//SetSAMOpts sets the SAM options
func SetSAMOpts(s []string) func(*Manager) error {
	return func(c *Manager) error {
		for _, i := range s {
			if i != "" {
				c.samopts = append(c.samopts, i)
			}
		}
		return nil
	}
}
