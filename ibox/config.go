package ibox

import (
	"fmt"
)

type Config struct {
	Username string
	Password string
	Hostname string
}

func (c *Config) Client() (*Client, error) {
	client, err := NewClient(c.Username, c.Password, c.Hostname)

	if err != nil {
		return nil, fmt.Errorf("[ERROR] setting up client failed: %s", err)
	}

	fmt.Printf("[INFO] Client configured for server %s", c.Hostname)

	return client, nil
}
