package client

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client connects to one or more Server using HTTP websockets.
// The Server can then send HTTP requests to execute.
type Client struct {
	Config *Config

	client *http.Client
	dialer *websocket.Dialer
	pools  map[string]*Pool
}

// NewClient creates a new Client.
func NewClient(config *Config) (c *Client) {
	c = new(Client)
	c.Config = config
	c.client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	c.dialer = &websocket.Dialer{}
	c.pools = make(map[string]*Pool)
	return
}

// Start the Proxy
func (c *Client) Start(ctx context.Context) {
	pool := NewPool(c, c.Config.Target)
	c.pools[c.Config.id] = pool
	go pool.Start(ctx)
}

// Shutdown the Proxy
func (c *Client) Shutdown() {
	for _, pool := range c.pools {
		pool.Shutdown()
	}
}
