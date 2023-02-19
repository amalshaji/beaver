package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

func getRequestHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

type PostPayload struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func postRequestHandler(c echo.Context) error {
	var postPayload PostPayload
	err := c.Bind(&postPayload)
	if err != nil {
		return c.JSON(400, map[string]string{"message": "invalid payload"})
	}
	return c.JSON(200, map[string]string{"message": "ok"})
}

func putRequestHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func patchRequestHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func deleteRequestHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func headRequestHandler(c echo.Context) error {
	c.Response().Header().Set("custom-server", "beaver-server")
	return c.NoContent(200)
}

func connectRequestHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func optionsRequestHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func redirect302RequestHandler(c echo.Context) error {
	return c.Redirect(302, "/")
}

func redirect307RequestHandler(c echo.Context) error {
	return c.Redirect(307, "/")
}

func main() {
	app := echo.New()
	app.HideBanner = true

	app.GET("/", getRequestHandler)
	app.POST("/", postRequestHandler)
	app.PUT("/", putRequestHandler)
	app.PATCH("/", patchRequestHandler)
	app.DELETE("/", deleteRequestHandler)
	app.OPTIONS("/", optionsRequestHandler)
	app.HEAD("/", headRequestHandler)
	app.CONNECT("/", connectRequestHandler)

	app.GET("/redirect-302", redirect302RequestHandler)
	app.GET("/redirect-307", redirect307RequestHandler)

	log.Fatal(app.Start(":9999"))
}
