package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var ErrAuthRequired = errors.New("authentication required")

type LoginPayload struct {
	Email    string `jspn:"email"`
	Password string `jspn:"password"`
}

func AuthRequired(c echo.Context) error {
	sessionToken, err := c.Request().Cookie("beaver_session")
	if err != nil {
		return ErrAuthRequired
	}
	if sessionToken == nil {
		return ErrAuthRequired
	}
	app := c.Get("app").(*App)
	_, err = app.User.ValidateSession(sessionToken.Value)
	if err != nil {
		return ErrAuthRequired
	}
	return nil
}

func authRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := AuthRequired(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return next(c)
	}
}

func SetupApiRoutes(e *echo.Echo) {
	g := e.Group("/api/v1")

	g.POST("/login", LoginApi)
	g.POST("/logout", LogoutApi, authRequiredMiddleware)
	g.GET("/stats", ServerStats, authRequiredMiddleware)
}

func LoginApi(c echo.Context) error {
	var p LoginPayload

	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	app := c.Get("app").(*App)

	token, err := app.User.Login(p.Email, p.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cookie := new(http.Cookie)
	cookie.Name = "beaver_session"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"

	c.SetCookie(cookie)

	c.JSON(200, map[string]string{"message": "ok"})

	return nil
}

func LogoutApi(c echo.Context) error {
	sessionCookie, err := c.Request().Cookie("beaver_session")

	if err != nil || sessionCookie == nil || sessionCookie.Value == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session"})
	}

	app := c.Get("app").(*App)

	err = app.User.Logout(sessionCookie.Value)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	sessionCookie.Value = ""
	sessionCookie.Expires = time.Now()
	sessionCookie.Path = "/"

	c.SetCookie(sessionCookie)

	return nil
}

func ServerStats(c echo.Context) error {
	var result = make(map[string]any)

	v, _ := mem.VirtualMemory()
	result["memory_used"] = v.UsedPercent

	_c, _ := cpu.Percent(0, false)
	result["cpu_used"] = _c

	app := c.Get("app").(*App)
	app.Server.Lock.Lock()
	defer app.Server.Lock.Unlock()

	result["active_connections"] = len(app.Server.Pools)

	c.JSON(200, result)
	return nil
}
