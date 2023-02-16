package handler

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/amalshaji/beaver/internal/server/app"
	"github.com/labstack/echo/v4"
)

func Start(configFile string) {
	e := echo.New()
	e.HideBanner = true

	_app := app.NewApp(configFile)

	adminHandler := GetAdminHandler(_app)
	tunnelHandler := GetTunnelHandler(_app)

	e.Any("/*", func(c echo.Context) error {
		req := c.Request()
		res := c.Response()

		_, err := _app.Server.GetSubdomainFromHost(req.Host)
		if err != nil {
			adminHandler.ServeHTTP(res, req)
		} else {
			tunnelHandler.ServeHTTP(res, req)
		}
		return nil
	})

	go func() { log.Fatal(e.Start(_app.Server.Config.GetAddr())) }()

	// Start the app
	_app.Start()

	// Wait signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	_app.Shutdown()
}
