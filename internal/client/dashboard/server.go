package dashboard

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/amalshaji/beaver/internal/utils"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/color"
)

var upgrader = websocket.Upgrader{}

func StartServer(connectionLogger *ConnectionLogger) *echo.Echo {
	e := echo.New()
	e.Logger.SetOutput(io.Discard)

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("dashboard", connectionLogger)
			return next(c)
		}
	})

	e.GET("/requests", func(c echo.Context) error {
		subdomain := c.QueryParam("subdomain")

		if err := connectionLogger.CheckSubdomain(subdomain); err != nil {
			return utils.HttpBadRequest(c, err.Error())
		}

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return utils.HttpBadRequest(c, err.Error())
		}

		defer func() {
			ws.Close()
			// remove the connection from the connection pool
		}()

		connectionLogger.AddConnection(subdomain, ws)

		for {
			time.Sleep(30 * time.Second)
			err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second))
			if err != nil {
				ws.Close()
			}
		}
	})

	go func() {
		log.Fatal(e.Start(":7878"))
	}()
	log.Println(color.Yellow("Dashboard running on http://localhost:7878"))

	return e
}

func StopServer(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
		}
	}
}
