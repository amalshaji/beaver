package server

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/amalshaji/beaver/internal/server/app"
	"github.com/amalshaji/beaver/internal/server/tunnel"
	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Request(c echo.Context) error {
	// [1]: Receive requests for proxying
	app := c.Get("app").(*app.App)

	subdomain, err := app.Server.GetSubdomainFromHost(c.Request().Host)

	if subdomain == "" {
		return utils.ProxyErrorf(c, "Subdomain required")
	}

	if err != nil {
		return utils.ProxyErrorf(c, err.Error())
	}

	// Parse destination URL
	var dstURL string

	if dstURL = app.Server.GetDestinationURL(subdomain); dstURL == "" {
		return utils.ProxyErrorf(c, "unregistered tunnel subdomain")
	}

	dstURL = fmt.Sprintf("%s/%s", dstURL, c.Param("*"))

	if c.QueryString() != "" {
		dstURL = fmt.Sprintf("%s?%s", dstURL, c.QueryString())
	}

	URL, err := url.Parse(dstURL)
	if err != nil {
		return utils.ProxyErrorf(c, "Unable to parse destination local server URL")
	}

	c.Request().URL = URL

	log.Printf("[%s] %s", c.Request().Method, c.Request().URL.String())

	if len(app.Server.Pools) == 0 {
		return utils.ProxyErrorf(c, "No proxy available")
	}

	// [2]: Take an WebSocket connection available from pools for relaying received requests.
	request := tunnel.NewConnectionRequest(app.Server.Config.GetTimeout())
	// "Dispatcher" is running in a separate thread from the server by `go s.dispatchConnections()`.
	// It waits to receive requests to dispatch connection from available pools to clients requests.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L93
	//
	// Notify request from handler to dispatcher through Server.dispatcher channel.
	app.Server.Dispatcher <- request
	// Dispatcher tries to find an available connection pool,
	// and it returns the connection through Server.connection channel.
	// https://github.com/hgsgtk/wsp/blob/ea4902a8e11f820268e52a6245092728efeffd7f/server/server.go#L189
	//
	// Here waiting for a result from dispatcher.

	connection := <-request.Connection
	if connection == nil {
		// It means that dispatcher has set `nil` which is a system error case that is
		// not expected in the normal flow.
		return utils.ProxyErrorf(c, "Unable to get a proxy connection")
	}

	// [3]: Send the request to the peer through the WebSocket connection.
	if err := connection.ProxyRequest(c); err != nil {
		// An error occurred throw the connection away
		log.Println(err)
		connection.Close()

		// Try to return an error to the client
		// This might fail if response headers have already been sent
		return utils.ProxyError(c, err)
	}
	return nil
}

// Request receives the WebSocket upgrade handshake request from wsp_client.
func Register(c echo.Context) error {
	app := c.Get("app").(*app.App)

	subdomain := c.Request().Header.Get("X-TUNNEL-SUBDOMAIN")
	localServer := c.Request().Header.Get("X-LOCAL-SERVER")

	secretKey := c.Request().Header.Get("X-SECRET-KEY")
	greeting := c.Request().Header.Get("X-GREETING-MESSAGE")

	var userIdentifier string
	for _, user := range app.Server.Config.Users {
		if user.SecretKey == secretKey {
			userIdentifier = user.Identifier
		}
	}

	if userIdentifier == "" {
		return utils.ProxyErrorf(c, "Invalid X-SECRET-KEY")
	}

	// Parse the greeting message
	split := strings.Split(string(greeting), "_")
	id := tunnel.PoolID(split[0])
	size, err := strconv.Atoi(split[1])
	if err != nil {
		return utils.ProxyErrorf(c, "Unable to parse greeting message : %s", err)
	}

	// 3. Register the connection into server pools.
	// s.lock is for exclusive control of pools operation.
	app.Server.Lock.Lock()
	defer app.Server.Lock.Unlock()

	pool, err := app.Server.GetOrCreatePoolForUser(subdomain, localServer, userIdentifier, id)
	if err != nil {
		return utils.ProxyErrorf(c, "subdomain already in use")
	}

	// update pool size
	pool.SetSize(size)

	// Upgrade the received HTTP request to a WebSocket connection
	ws, err := app.Server.Upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return utils.ProxyErrorf(c, "HTTP upgrade error : %v", err)
	}

	// Add the WebSocket connection to the pool
	pool.Register(ws)

	return nil
}

func status(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func Start(configFile string) {
	e := echo.New()
	e.HideBanner = true

	app := app.NewApp(configFile)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	e.GET("/register", Register)
	e.GET("/status", status)

	// Handle tunnel requests
	e.GET("*", Request)
	e.POST("*", Request)
	e.PUT("*", Request)
	e.PATCH("*", Request)
	e.DELETE("*", Request)
	e.OPTIONS("*", Request)

	e.Use(middleware.Recover())

	// Dispatch connection from available pools to clients requests
	// in a separate thread from the server thread.
	go app.Server.DispatchConnections()

	go func() { log.Fatal(e.Start(app.Server.Config.GetAddr())) }()

	// Wait signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	app.Shutdown()
}
