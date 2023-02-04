package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"github.com/amalshaji/beaver"
)

// Status of a Connection
const (
	CONNECTING = iota
	IDLE
	RUNNING
)

// Connection handle a single websocket (HTTP/TCP) connection to an Server
type Connection struct {
	pool   *Pool
	ws     *websocket.Conn
	status int
}

func (c *Connection) IsInitialConnection() bool {
	return len(c.pool.connections) == 1
}

// NewConnection create a Connection object
func NewConnection(pool *Pool) *Connection {
	c := new(Connection)
	c.pool = pool
	c.status = CONNECTING
	return c
}

// Connect to the IsolatorServer using a HTTP websocket
func (connection *Connection) Connect(ctx context.Context) (err error) {
	if connection.IsInitialConnection() {
		log.Println("Creating tunnel connection")
	}

	var res *http.Response
	// Create a new TCP(/TLS) connection ( no use of net.http )
	connection.ws, res, err = connection.pool.client.dialer.DialContext(
		ctx,
		connection.pool.target,
		http.Header{
			"X-SECRET-KEY":       {connection.pool.client.Config.SecretKey},
			"X-TUNNEL-SUBDOMAIN": {connection.pool.client.Config.subdomain},
			"X-LOCAL-SERVER": {fmt.Sprintf(
				"http://localhost:%d",
				connection.pool.client.Config.port,
			)},
			"X-GREETING-MESSAGE": {fmt.Sprintf(
				"%s_%d",
				connection.pool.client.Config.id,
				connection.pool.client.Config.PoolIdleSize,
			)},
		},
	)

	if err != nil {
		defer func() {
			if r := recover(); r != nil {
				log.Fatal(err)
			}
		}()
		bodyBytes, _ := io.ReadAll(res.Body)
		defer res.Body.Close()
		log.Fatal(string(bodyBytes))
	}

	var httpScheme string

	URL, _ := url.Parse(connection.pool.target)
	if URL.Scheme == "ws" {
		httpScheme = "http"
	} else {
		httpScheme = "https"
	}
	httpPort := URL.Port()
	if httpPort != "" {
		httpPort = ":" + httpPort
	}

	if connection.IsInitialConnection() {
		log.Printf("Tunnel running on %s://%s.%s%s",
			httpScheme,
			connection.pool.client.Config.subdomain,
			URL.Hostname(),
			httpPort,
		)
	}

	// Send the greeting message with proxy id and wanted pool size.

	go connection.serve(ctx)

	return
}

// the main loop it :
//   - wait to receive HTTP requests from the Server
//   - execute HTTP requests
//   - send HTTP response back to the Server
//
// As in the server code there is no buffering of HTTP request/response body
// As is the server if any error occurs the connection is closed/throwed
func (connection *Connection) serve(ctx context.Context) {
	defer connection.Close()

	// Keep connection alive
	go func() {
		for {
			time.Sleep(30 * time.Second)
			err := connection.ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				connection.Close()
			}
		}
	}()

	for {
		// Read request
		connection.status = IDLE
		_, jsonRequest, err := connection.ws.ReadMessage()
		if err != nil {
			if connection.pool.client.Config.showWsReadErrors {
				log.Println("Unable to read request", err)
			}
			break
		}

		connection.status = RUNNING

		// Trigger a pool refresh to open new connections if needed
		go connection.pool.connector(ctx)

		// Deserialize request
		httpRequest := new(beaver.HTTPRequest)
		err = json.Unmarshal(jsonRequest, httpRequest)
		if err != nil {
			connection.error(fmt.Sprintf("Unable to deserialize json http request : %s\n", err))
			break
		}

		req, err := beaver.UnserializeHTTPRequest(httpRequest)
		if err != nil {
			connection.error(fmt.Sprintf("Unable to deserialize http request : %v\n", err))
			break
		}

		// Pipe request body
		_, bodyReader, err := connection.ws.NextReader()
		if err != nil {
			log.Printf("Unable to get response body reader : %v", err)
			break
		}
		req.Body = io.NopCloser(bodyReader)

		// Execute request
		resp, err := connection.pool.client.client.Do(req)
		if err != nil {
			err = connection.error(fmt.Sprintf("Unable to execute request : %v\n", err))
			if err != nil {
				break
			}
			continue
		}

		var urlPath string = req.URL.Path

		if req.URL.RawQuery != "" {
			urlPath = urlPath + "?" + req.URL.RawQuery
		}

		log.Printf("[%s] %d %s", req.Method, resp.StatusCode, urlPath)

		// Serialize response
		jsonResponse, err := json.Marshal(beaver.SerializeHTTPResponse(resp))
		if err != nil {
			err = connection.error(fmt.Sprintf("Unable to serialize response : %v\n", err))
			if err != nil {
				break
			}
			continue
		}

		// Write response
		err = connection.ws.WriteMessage(websocket.TextMessage, jsonResponse)
		if err != nil {
			log.Printf("Unable to write response : %v", err)
			break
		}

		// Pipe response body
		bodyWriter, err := connection.ws.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Printf("Unable to get response body writer : %v", err)
			break
		}
		_, err = io.Copy(bodyWriter, resp.Body)
		if err != nil {
			log.Printf("Unable to get pipe response body : %v", err)
			break
		}
		bodyWriter.Close()
	}
}

func (connection *Connection) error(msg string) (err error) {
	resp := beaver.NewHTTPResponse()
	resp.StatusCode = 527

	log.Println(msg)

	resp.ContentLength = int64(len(msg))

	// Serialize response
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Unable to serialize response : %v", err)
		return
	}

	// Write response
	err = connection.ws.WriteMessage(websocket.TextMessage, jsonResponse)
	if err != nil {
		log.Printf("Unable to write response : %v", err)
		return
	}

	// Write response body
	err = connection.ws.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		log.Printf("Unable to write response body : %v", err)
		return
	}

	return
}

// Close close the ws/tcp connection and remove it from the pool
func (connection *Connection) Close() {
	connection.pool.lock.Lock()
	defer connection.pool.lock.Unlock()

	connection.pool.remove(connection)
	connection.ws.Close()
}
