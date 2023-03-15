package dashboard

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ErrConnectionExists       = errors.New("connection exists")
	ErrConnectionDoesNotExist = errors.New("connection does not exist")
	MaxLogsCount              = 100
)

type ConnectionLogger struct {
	connections map[string][]*websocket.Conn
	lock        sync.Mutex
}

func NewConnectionLogger() *ConnectionLogger {
	return &ConnectionLogger{
		connections: make(map[string][]*websocket.Conn),
	}
}

func (c *ConnectionLogger) AddConnection(subdomain string, connection *websocket.Conn) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.CheckSubdomain(subdomain)
	if err != nil {
		c.connections[subdomain] = make([]*websocket.Conn, 0)
	}

	if connection != nil {
		c.connections[subdomain] = append(c.connections[subdomain], connection)
	}

	return nil
}

func (c *ConnectionLogger) CheckSubdomain(subdomain string) error {
	_, ok := c.connections[subdomain]
	if !ok {
		return ErrConnectionDoesNotExist
	}
	return nil
}

func (c *ConnectionLogger) IsNewConnection(subdomain string) bool {
	_, ok := c.connections[subdomain]
	return !ok
}

func (c *ConnectionLogger) LogRequest(subdomain string, request *http.Request, response *http.Response) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	connections, ok := c.connections[subdomain]
	if !ok {
		return ErrConnectionDoesNotExist
	}

	requestBytes, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}

	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		return err
	}

	newRequest := map[string]any{
		"request":  string(requestBytes),
		"response": string(responseBytes),
	}

	for _, connection := range connections {
		err := connection.WriteJSON(newRequest)
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
