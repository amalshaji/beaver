package handler

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amalshaji/beaver/internal/server/app"
	"github.com/amalshaji/beaver/internal/server/tunnel"
	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/echo/v4"
)

func request(c echo.Context) error {
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
	request := tunnel.NewConnectionRequest(app.Server.Config.GetTimeout(), subdomain)
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

func GetTunnelHandler(app *app.App) *echo.Echo {
	tunnelRouter := echo.New()

	tunnelRouter.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	tunnelRouter.GET("*", request)
	tunnelRouter.POST("*", request)
	tunnelRouter.PUT("*", request)
	tunnelRouter.PATCH("*", request)
	tunnelRouter.DELETE("*", request)
	tunnelRouter.OPTIONS("*", request)
	tunnelRouter.HEAD("*", request)
	tunnelRouter.CONNECT("*", request)

	return tunnelRouter
}
