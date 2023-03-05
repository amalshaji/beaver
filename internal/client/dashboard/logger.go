package dashboard

import (
	"container/list"
	"errors"
	"net/http"
	"net/http/httputil"
	"sync"
)

var (
	ErrConnectionExists       = errors.New("connection exists")
	ErrConnectionDoesNotExist = errors.New("connection does not exist")
	MaxLogsCount              = 100
)

type ConnectionLogger struct {
	connections map[string]*requestLogs
}

type requestLogs struct {
	logs *list.List
	lock sync.Mutex
}

func NewConnectionLogger() *ConnectionLogger {
	return &ConnectionLogger{
		connections: make(map[string]*requestLogs),
	}
}

func (c *ConnectionLogger) AddConnection(subdomain string) error {
	_, ok := c.connections[subdomain]
	if ok {
		return ErrConnectionExists
	}

	c.connections[subdomain] = &requestLogs{
		logs: list.New(),
	}
	return nil
}

func (c *ConnectionLogger) IsNewConnection(subdomain string) bool {
	_, ok := c.connections[subdomain]
	return !ok
}

func (c *ConnectionLogger) LogRequest(subdomain string, request *http.Request, response *http.Response) error {
	connection, ok := c.connections[subdomain]
	if !ok {
		return ErrConnectionDoesNotExist
	}

	connection.lock.Lock()
	defer connection.lock.Unlock()

	// Remove first inserted(last pos) connection if max limit is reached
	if connection.logs.Len() >= MaxLogsCount {
		front := connection.logs.Back()
		connection.logs.Remove(front)
	}

	requestBytes, err := httputil.DumpRequest(request, true)
	if err != nil {
		return err
	}

	responseBytes, err := httputil.DumpResponse(response, true)
	if err != nil {
		return err
	}

	connection.logs.PushFront(map[string][]byte{
		"request":  requestBytes,
		"response": responseBytes,
	})

	return nil
}

func (c *ConnectionLogger) GetLogs(subdomain string) (*list.List, error) {
	connection, ok := c.connections[subdomain]
	if !ok {
		return nil, ErrConnectionDoesNotExist
	}

	connection.lock.Lock()
	defer connection.lock.Unlock()

	return connection.logs, nil
}
