package tunnel

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Server is a Reverse HTTP Proxy over WebSocket
// This is the Server part, Clients will offer websocket connections,
// those will be pooled to transfer HTTP Request and response
type Server struct {
	Config *Config

	Upgrader websocket.Upgrader

	// In Pools, keep connections with WebSocket peers.
	Pools map[string]*Pool

	// A RWMutex is a reader/writer mutual exclusion Lock,
	// and it is for exclusive control with pools operation.
	//
	// This is locked when reading and writing pools, the timing is when:
	// 1. (rw) registering websocket clients in /register endpoint
	// 2. (rw) remove empty pools which has no connections
	// 3. (r) dispatching connection from available pools to clients requests
	//
	// And then it is released after each process is completed.
	Lock sync.RWMutex
	done chan struct{}

	// Through Dispatcher channel it communicates between "server" thread and "Dispatcher" thread.
	// "server" thread sends the value to this channel when accepting requests in the endpoint /requests,
	// and "Dispatcher" thread reads this channel.
	Dispatcher chan *ConnectionRequest
}

// ConnectionRequest is used to request a proxy connection from the dispatcher
type ConnectionRequest struct {
	// Pick the connection for this subdomain
	Subdomain string

	Connection chan *Connection
}

// NewConnectionRequest creates a new connection request
func NewConnectionRequest(timeout time.Duration, connectionFor string) (cr *ConnectionRequest) {
	cr = new(ConnectionRequest)
	cr.Connection = make(chan *Connection)
	cr.Subdomain = connectionFor
	return
}

// NewServer return a new Server instance
func NewServer(configFile string) (server *Server) {
	rand.Seed(time.Now().Unix())

	// Load configuration
	config, err := LoadConfiguration(configFile)
	if err != nil {
		log.Fatalf("Unable to load configuration : %s", err)
	}

	server = new(Server)
	server.Config = config
	server.Upgrader = websocket.Upgrader{}
	server.Pools = make(map[string]*Pool)

	server.done = make(chan struct{})
	server.Dispatcher = make(chan *ConnectionRequest)

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

	// Dispatch connection from available pools to clients requests
	// in a separate thread from the server thread.
	go s.DispatchConnections()
}

// clean removes empty Pools which has no connection.
// It is invoked every 5 sesconds and at shutdown.
func (s *Server) clean() {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if len(s.Pools) == 0 {
		return
	}

	idle := 0
	busy := 0

	pools := make(map[string]*Pool)
	for subdomain, pool := range s.Pools {
		if pool.IsEmpty() {
			log.Printf("Removing empty connection pool : %s", pool.ID)
			pool.Shutdown()
		} else {
			pools[subdomain] = pool
		}

		ps := pool.Size()
		idle += ps.Idle
		busy += ps.Busy
	}

	log.Printf("%d pools, %d idle, %d busy", len(pools), idle, busy)

	s.Pools = pools
}

// Dispatch connection from available pools to clients requests
func (s *Server) DispatchConnections() {
	for {
		// Runs in an infinite loop and keeps receiving the value from the `server.dispatcher` channel
		// The operator <- is "receive operator", which expression blocks until a value is available.
		request, ok := <-s.Dispatcher
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

			s.Lock.RLock()
			pool, ok := s.Pools[request.Subdomain]
			if !ok {
				// No connection pool available
				s.Lock.RUnlock()
				break
			}

			var connection *Connection
			for _, conn := range pool.connections {
				if conn.status == Idle {
					connection = conn
				}
			}
			s.Lock.RUnlock()

			if connection == nil {
				continue
			}

			// [2]: Verify that we can use this connection and take it.
			if connection.Take() {
				request.Connection <- connection
				break
			}
		}

		close(request.Connection)
	}
}

func (s *Server) GetSubdomainFromHost(host string) (string, error) {
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

// Shutdown stop the Server
func (s *Server) Shutdown() {
	close(s.done)
	close(s.Dispatcher)
	for _, pool := range s.Pools {
		pool.Shutdown()
	}
	s.clean()
}

func (s *Server) GetOrCreatePoolForUser(subdomain, localServer, userIdentifier string, id PoolID) (*Pool, error) {
	// There is no need to create a new pool,
	// if it is already registered in current pools.
	p, ok := s.Pools[subdomain]
	if !ok {
		pool := NewPool(s, id, subdomain, localServer, userIdentifier)
		s.Pools[subdomain] = pool
		return pool, nil
	}

	if p.ID != id {
		return nil, fmt.Errorf("subdomain already in use")
	}

	return p, nil
}

func (s *Server) GetDestinationURL(subdomain string) string {
	p, ok := s.Pools[subdomain]
	if !ok {
		return ""
	}

	return p.LocalServer
}
