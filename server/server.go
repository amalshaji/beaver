package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	beaver "github.com/amalshaji/beaver"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

// Server is a Reverse HTTP Proxy over WebSocket
// This is the Server part, Clients will offer websocket connections,
// those will be pooled to transfer HTTP Request and response
type Server struct {
	Config *Config

	upgrader websocket.Upgrader

	// In pools, keep connections with WebSocket peers.
	pools []*Pool

	// A RWMutex is a reader/writer mutual exclusion lock,
	// and it is for exclusive control with pools operation.
	//
	// This is locked when reading and writing pools, the timing is when:
	// 1. (rw) registering websocket clients in /register endpoint
	// 2. (rw) remove empty pools which has no connections
	// 3. (r) dispatching connection from available pools to clients requests
	//
	// And then it is released after each process is completed.
	lock sync.RWMutex
	done chan struct{}

	// Through dispatcher channel it communicates between "server" thread and "dispatcher" thread.
	// "server" thread sends the value to this channel when accepting requests in the endpoint /requests,
	// and "dispatcher" thread reads this channel.
	dispatcher chan *ConnectionRequest
}

// ConnectionRequest is used to request a proxy connection from the dispatcher
type ConnectionRequest struct {
	connection chan *Connection
}

// NewConnectionRequest creates a new connection request
func NewConnectionRequest(timeout time.Duration) (cr *ConnectionRequest) {
	cr = new(ConnectionRequest)
	cr.connection = make(chan *Connection)
	return
}

// NewServer return a new Server instance
func NewServer(config *Config) (server *Server) {
	rand.Seed(time.Now().Unix())

	server = new(Server)
	server.Config = config
	server.upgrader = websocket.Upgrader{}

	server.done = make(chan struct{})
	server.dispatcher = make(chan *ConnectionRequest)
	return
}

// Start Server HTTP server
func (s *Server) Start() {
	go func() {
	L:
		for {
			select {
			case <-s.done:
				break L
			case <-time.After(5 * time.Second):
				s.clean()
			}
		}
	}()

	e := echo.New()
	e.HideBanner = true

	e.GET("/register", s.Register)
	e.GET("/status", s.status)

	// Handle tunnel requests
	e.GET("*", s.Request)
	e.POST("*", s.Request)
	e.PUT("*", s.Request)
	e.PATCH("*", s.Request)
	e.DELETE("*", s.Request)
	e.OPTIONS("*", s.Request)

	// Dispatch connection from available pools to clients requests
	// in a separate thread from the server thread.
	go s.dispatchConnections()

	go func() { log.Fatal(e.Start(s.Config.GetAddr())) }()
}

// clean removes empty Pools which has no connection.
// It is invoked every 5 sesconds and at shutdown.
func (s *Server) clean() {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.pools) == 0 {
		return
	}

	idle := 0
	busy := 0

	var pools []*Pool
	for _, pool := range s.pools {
		if pool.IsEmpty() {
			log.Printf("Removing empty connection pool : %s", pool.id)
			pool.Shutdown()
		} else {
			pools = append(pools, pool)
		}

		ps := pool.Size()
		idle += ps.Idle
		busy += ps.Busy
	}

	log.Printf("%d pools, %d idle, %d busy", len(pools), idle, busy)

	s.pools = pools
}

// Dispatch connection from available pools to clients requests
func (s *Server) dispatchConnections() {
	for {
		// Runs in an infinite loop and keeps receiving the value from the `server.dispatcher` channel
		// The operator <- is "receive operator", which expression blocks until a value is available.
		request, ok := <-s.dispatcher
		if !ok {
			// The value of `ok` is false if it is a zero value generated because the channel is closed an empty.
			// In this case, that means server shutdowns.
			break
		}

		// A timeout is set for each dispatch request.
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, s.Config.GetTimeout())
		defer cancel()

	L:
		for {
			select {
			case <-ctx.Done(): // The timeout elapses
				break L
			default: // Go through
			}

			s.lock.RLock()
			if len(s.pools) == 0 {
				// No connection pool available
				s.lock.RUnlock()
				break
			}

			// [1]: Select a pool which has an idle connection
			// Build a select statement dynamically to handle an arbitrary number of pools.
			cases := make([]reflect.SelectCase, len(s.pools)+1)
			for i, ch := range s.pools {
				cases[i] = reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(ch.idle)}
			}
			cases[len(cases)-1] = reflect.SelectCase{
				Dir: reflect.SelectDefault}
			s.lock.RUnlock()

			_, value, ok := reflect.Select(cases)
			if !ok {
				continue // a pool has been removed, try again
			}
			connection, _ := value.Interface().(*Connection)

			// [2]: Verify that we can use this connection and take it.
			if connection.Take() {
				request.connection <- connection
				break
			}
		}

		close(request.connection)
	}
}

func (s *Server) getSubdomainFromHost(host string) (string, error) {
	var httpScheme string
	if s.Config.Secure {
		httpScheme = "https"
	} else {
		httpScheme = "http"
	}
	if !strings.HasPrefix(host, httpScheme+"://") {
		host = fmt.Sprintf("%s://%s", httpScheme, host)
	}
	url, err := url.Parse(host)

	if err != nil {
		return "", err
	}

	hostname := url.Hostname()

	if !strings.HasSuffix(hostname, "."+s.Config.Domain) {
		return "", fmt.Errorf("subdomain required")
	}
	return strings.Replace(hostname, "."+s.Config.Domain, "", 1), nil
}

func (s *Server) Request(c echo.Context) error {
	// [1]: Receive requests to be proxied
	// Parse destination URL
	var dstURL string

	subdomain, err := s.getSubdomainFromHost(c.Request().Host)

	if err != nil {
		return beaver.ProxyErrorf(c, err.Error())
	}

	for _, p := range s.pools {
		if p.subdomain == subdomain {
			dstURL = p.localServer
			break
		}
	}

	if dstURL == "" {
		return beaver.ProxyErrorf(c, "unregistered tunnel subdomain")
	}

	dstURL = fmt.Sprintf("%s/%s", dstURL, c.Param("*"))

	if c.QueryString() != "" {
		dstURL = fmt.Sprintf("%s?%s", dstURL, c.QueryString())
	}
	if dstURL == "" {
		return beaver.ProxyErrorf(c, "Subdomain required")
	}
	URL, err := url.Parse(dstURL)
	if err != nil {
		return beaver.ProxyErrorf(c, "Unable to parse destination local server URL")
	}

	c.Request().URL = URL

	log.Printf("[%s] %s", c.Request().Method, c.Request().URL.String())

	if len(s.pools) == 0 {
		return beaver.ProxyErrorf(c, "No proxy available")
	}

	// [2]: Take an WebSocket connection available from pools for relaying received requests.
	request := NewConnectionRequest(s.Config.GetTimeout())
	// "Dispatcher" is running in a separate thread from the server by `go s.dispatchConnections()`.
	// It waits to receive requests to dispatch connection from available pools to clients requests.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L93
	//
	// Notify request from handler to dispatcher through Server.dispatcher channel.
	s.dispatcher <- request
	// Dispatcher tries to find an available connection pool,
	// and it returns the connection through Server.connection channel.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L189
	//
	// Here waiting for a result from dispatcher.
	connection := <-request.connection
	if connection == nil {
		// It means that dispatcher has set `nil` which is a system error case that is
		// not expected in the normal flow.
		return beaver.ProxyErrorf(c, "Unable to get a proxy connection")
	}

	// [3]: Send the request to the peer through the WebSocket connection.
	if err := connection.proxyRequest(c); err != nil {
		// An error occurred throw the connection away
		log.Println(err)
		connection.Close()

		// Try to return an error to the client
		// This might fail if response headers have already been sent
		return beaver.ProxyError(c, err)
	}
	return nil
}

// Request receives the WebSocket upgrade handshake request from wsp_client.
func (s *Server) Register(c echo.Context) error {
	// 1. Upgrade a received HTTP request to a WebSocket connection
	subdomain := c.Request().Header.Get("X-TUNNEL-SUBDOMAIN")
	localServer := c.Request().Header.Get("X-LOCAL-SERVER")

	secretKey := c.Request().Header.Get("X-SECRET-KEY")
	greeting := c.Request().Header.Get("X-GREETING-MESSAGE")

	var userIdentifier string
	for _, user := range s.Config.Users {
		if user.SecretKey == secretKey {
			userIdentifier = user.Identifier
		}
	}

	if userIdentifier == "" {
		return beaver.ProxyErrorf(c, "Invalid X-SECRET-KEY")
	}

	// Parse the greeting message
	split := strings.Split(string(greeting), "_")
	id := PoolID(split[0])
	size, err := strconv.Atoi(split[1])
	if err != nil {
		return beaver.ProxyErrorf(c, "Unable to parse greeting message : %s", err)
	}

	// 3. Register the connection into server pools.
	// s.lock is for exclusive control of pools operation.
	s.lock.Lock()
	defer s.lock.Unlock()

	var pool *Pool

	// There is no need to create a new pool,
	// if it is already registered in current pools.
	for _, p := range s.pools {
		if p.subdomain == subdomain {
			if p.id == id {
				pool = p
				break
			} else {
				return beaver.ProxyErrorf(c, "subdomain already in use")
			}
		}
	}
	if pool == nil {
		pool = NewPool(s, id, subdomain, localServer, userIdentifier)
		s.pools = append(s.pools, pool)
	}
	// update pool size
	pool.size = size

	ws, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return beaver.ProxyErrorf(c, "HTTP upgrade error : %v", err)
	}

	// Add the WebSocket connection to the pool
	pool.Register(ws)

	return nil
}

func (s *Server) status(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

// Shutdown stop the Server
func (s *Server) Shutdown() {
	close(s.done)
	close(s.dispatcher)
	for _, pool := range s.pools {
		pool.Shutdown()
	}
	s.clean()
}
