package dashboard

import (
	"context"
	"io"
	"log"
	"strconv"
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
		lastIdStr := c.QueryParam("lastId")
		if lastIdStr == "" {
			lastIdStr = "0"
		}
		lastId, err := strconv.Atoi(lastIdStr)
		if err != nil {
			return utils.HttpBadRequest(c, "lastId must be a number")
		}

		logs, err := connectionLogger.GetLogs(c.QueryParam("subdomain"))
		if err != nil {
			return utils.HttpBadRequest(c, err.Error())
		}

		var s []map[string]any
		for e := logs.Front(); e != nil; e = e.Next() {
			id := e.Value.(map[string]any)["id"]
			if id.(int) <= lastId {
				continue
			}
			s = append(s, map[string]any{
				"id":       id,
				"request":  e.Value.(map[string]any)["request"],
				"response": e.Value.(map[string]any)["response"],
			})
		}

		if len(s) == 0 {
			return c.JSON(200, []string{})
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
