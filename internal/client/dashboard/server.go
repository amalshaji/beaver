package dashboard

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/echo/v4"
)

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
		logs, err := connectionLogger.GetLogs(c.QueryParam("subdomain"))
		if err != nil {
			return utils.HttpBadRequest(c, err.Error())
		}

		if logs.Len() == 0 {
			return c.JSON(200, []string{})
		}

		var s []map[string]string
		for e := logs.Front(); e != nil; e = e.Next() {
			s = append(s, map[string]string{
				"request":  e.Value.(map[string]string)["request"],
				"response": e.Value.(map[string]string)["response"],
			})
		}

		return c.JSON(200, s)
	})

	go func() {
		log.Fatal(e.Start(":7878"))
	}()
	log.Println("Dashboard running on http://localhost:7878")

	return e
}

func StopServer(e *echo.Echo) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
