package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/root-gg/wsp"
)

// ConnectionStatus is an enumeration type which represents the status of WebSocket connection.
type ConnectionStatus int

const (
	// Idle state means it is opened but not working now.
	// The default value for Connection is Idle, so it is ok to use zero-value(int: 0) for Idle status.
	Idle ConnectionStatus = iota
	Busy
	Closed
)

// Connection manages a single websocket connection from the peer.
// wsp supports multiple connections from a single peer at the same time.
type Connection struct {
	pool         *Pool
	ws           *websocket.Conn
	status       ConnectionStatus
	idleSince    time.Time
	lock         sync.Mutex
	nextResponse chan chan io.Reader
}

// NewConnection returns a new Connection.
func NewConnection(pool *Pool, ws *websocket.Conn) *Connection {
	// Initialize a new Connection
	c := new(Connection)
	c.pool = pool
	c.ws = ws
	c.nextResponse = make(chan chan io.Reader)
	c.status = Idle

	// Mark that this connection is ready to use for relay
	c.Release()

	// Start to listen to incoming messages over the WebSocket connection
	go c.read()

	return c
}

// read the incoming message of the connection
func (connection *Connection) read() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Websocket crash recovered : %s", r)
		}
		connection.Close()
	}()

	for {
		if connection.status == Closed {
			break
		}

		// https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
		//
		// We need to ensure :
		//  - no concurrent calls to ws.NextReader() / ws.ReadMessage()
		//  - only one reader exists at a time
		//  - wait for reader to be consumed before requesting the next one
		//  - always be reading on the socket to be able to process control messages ( ping / pong / close )

		// We will block here until a message is received or the ws is closed
		_, reader, err := connection.ws.NextReader()
		if err != nil {
			break
		}

		if connection.status != Busy {
			// We received a wild unexpected message
			break
		}

		// We received a message from the proxy
		// It is expected to be either a HttpResponse or a HttpResponseBody
		// We wait for proxyRequest to send a channel to get the message
		c := <-connection.nextResponse
		if c == nil {
			// We have been unlocked by Close()
			break
		}

		// Send the reader back to proxyRequest
		c <- reader

		// Wait for proxyRequest to close the channel
		// this notify that it is done with the reader
		<-c
	}
}

// Proxy a HTTP request through the Proxy over the websocket connection
func (connection *Connection) proxyRequest(w http.ResponseWriter, r *http.Request) (err error) {
	log.Printf("proxy request to %s", connection.pool.id)

	// [1]: Serialize HTTP request
	jsonReq, err := json.Marshal(wsp.SerializeHTTPRequest(r))
	if err != nil {
		return fmt.Errorf("unable to serialize request : %w", err)
	}
	// i.e.
	// {
	// 		"Method":"GET",
	// 		"URL":"http://localhost:8081/hello",
	// 		"Header":{"Accept":["*/*"],"User-Agent":["curl/7.77.0"],"X-Proxy-Destination":["http://localhost:8081/hello"]},
	//		"ContentLength":0
	// }

	// [2]: Send the HTTP request to the peer
	// Send the serialized HTTP request to the the peer
	if err := connection.ws.WriteMessage(websocket.TextMessage, jsonReq); err != nil {
		return fmt.Errorf("unable to write request : %w", err)
	}

	// Pipe the HTTP request body to the the peer
	bodyWriter, err := connection.ws.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return fmt.Errorf("unable to get request body writer : %w", err)
	}
	if _, err := io.Copy(bodyWriter, r.Body); err != nil {
		return fmt.Errorf("unable to pipe request body : %w", err)
	}
	if err := bodyWriter.Close(); err != nil {
		return fmt.Errorf("unable to pipe request body (close) : %w", err)
	}

	// [3]: Read the HTTP response from the peer
	// Get the serialized HTTP Response from the peer
	// To do so send a new channel to the read() goroutine
	// to get the next message reader
	responseChannel := make(chan (io.Reader))
	connection.nextResponse <- responseChannel
	responseReader, more := <-responseChannel
	if responseReader == nil {
		if more {
			// If more is false the channel is already closed
			close(responseChannel)
		}
		return fmt.Errorf("unable to get http response reader : %w", err)
	}

	// Read the HTTP Response
	jsonResponse, err := io.ReadAll(responseReader)
	if err != nil {
		close(responseChannel)
		return fmt.Errorf("unable to read http response : %w", err)
	}

	// Notify the read() goroutine that we are done reading the response
	close(responseChannel)

	// Deserialize the HTTP Response
	httpResponse := new(wsp.HTTPResponse)
	if err := json.Unmarshal(jsonResponse, httpResponse); err != nil {
		return fmt.Errorf("unable to unserialize http response : %w", err)
	}

	// Write response headers back to the client
	for header, values := range httpResponse.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}
	w.WriteHeader(httpResponse.StatusCode)

	// Get the HTTP Response body from the the peer
	// To do so send a new channel to the read() goroutine
	// to get the next message reader
	responseBodyChannel := make(chan (io.Reader))
	connection.nextResponse <- responseBodyChannel
	responseBodyReader, more := <-responseBodyChannel
	if responseBodyReader == nil {
		if more {
			// If more is false the channel is already closed
			close(responseChannel)
		}
		return fmt.Errorf("unable to get http response body reader : %w", err)
	}

	// Pipe the HTTP response body right from the remote Proxy to the client
	if _, err := io.Copy(w, responseBodyReader); err != nil {
		close(responseBodyChannel)
		return fmt.Errorf("unable to pipe response body : %w", err)
	}

	// Notify read() that we are done reading the response body
	close(responseBodyChannel)

	connection.Release()

	return
}

// Take notifies that this connection is going to be used
func (connection *Connection) Take() bool {
	connection.lock.Lock()
	defer connection.lock.Unlock()

	if connection.status == Closed {
		return false
	}

	if connection.status == Busy {
		return false
	}

	connection.status = Busy
	return true
}

// Release notifies that this connection is ready to use again
func (connection *Connection) Release() {
	connection.lock.Lock()
	defer connection.lock.Unlock()

	if connection.status == Closed {
		return
	}

	connection.idleSince = time.Now()
	connection.status = Idle

	go connection.pool.Offer(connection)
}

// Close the connection
func (connection *Connection) Close() {
	connection.lock.Lock()
	defer connection.lock.Unlock()

	connection.close()
}

// Close the connection ( without lock )
func (connection *Connection) close() {
	if connection.status == Closed {
		return
	}

	log.Printf("Closing connection from %s", connection.pool.id)

	// This one will be executed *before* lock.Unlock()
	defer func() { connection.status = Closed }()

	// Unlock a possible read() wild message
	close(connection.nextResponse)

	// Close the underlying TCP connection
	connection.ws.Close()
}
