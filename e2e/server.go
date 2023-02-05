package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

func getRequest(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

type PostPayload struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func postRequest(c echo.Context) error {
	var postPayload PostPayload
	err := c.Bind(&postPayload)
	if err != nil {
		return c.JSON(400, map[string]string{"message": "invalid payload"})
	}
	return c.JSON(200, map[string]string{"message": "ok"})
}

func putRequest(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func patchRequest(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func deleteRequest(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func optionsRequest(c echo.Context) error {
	return c.JSON(200, map[string]string{"message": "ok"})
}

func redirect302Request(c echo.Context) error {
	return c.Redirect(302, "/")
}

func redirect307Request(c echo.Context) error {
	return c.Redirect(307, "/")
}

func main() {
	app := echo.New()
	app.HideBanner = true

	app.GET("/", getRequest)
	app.POST("/", postRequest)
	app.PUT("/", putRequest)
	app.PATCH("/", patchRequest)
	app.DELETE("/", deleteRequest)
	app.OPTIONS("/", optionsRequest)

	app.GET("/redirect-302", redirect302Request)
	app.GET("/redirect-307", redirect307Request)

	log.Fatal(app.Start(":9999"))
}
