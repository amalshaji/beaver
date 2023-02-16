package handler

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/amalshaji/beaver/internal/server/admin"
	"github.com/amalshaji/beaver/internal/server/app"
	"github.com/amalshaji/beaver/internal/server/static"
	"github.com/amalshaji/beaver/internal/server/tunnel"
	"github.com/amalshaji/beaver/internal/server/web"
	"github.com/amalshaji/beaver/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Request receives the WebSocket upgrade handshake request from beaver client.
func register(c echo.Context) error {
	app := c.Get("app").(*app.App)

	subdomain := c.Request().Header.Get("X-TUNNEL-SUBDOMAIN")
	localServer := c.Request().Header.Get("X-LOCAL-SERVER")

	secretKey := c.Request().Header.Get("X-SECRET-KEY")
	greeting := c.Request().Header.Get("X-GREETING-MESSAGE")

	tunnelUser, err := app.User.GetTunnelUserBySecret(c.Request().Context(), secretKey)
	if err != nil && errors.Is(err, admin.ErrTunnelUserNotFound) {
		return utils.ProxyErrorf(c, "invalid secretKey - unregistered tunnel user")
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

	pool, err := app.Server.GetOrCreatePoolForUser(subdomain, localServer, tunnelUser.Email, id)
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

var ErrAuthRequired = errors.New("authentication required")

type LoginPayload struct {
	Email    string `jspn:"email"`
	Password string `jspn:"password"`
}

func authRequired(c echo.Context) error {
	sessionToken, err := c.Request().Cookie("beaver_session")
	if err != nil {
		return ErrAuthRequired
	}
	if sessionToken == nil {
		return ErrAuthRequired
	}
	app := c.Get("app").(*app.App)
	_, err = app.User.ValidateSession(c.Request().Context(), sessionToken.Value)
	if err != nil {
		return ErrAuthRequired
	}
	return nil
}

func authRequiredMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := authRequired(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		}
		return next(c)
	}
}

func setupApiRoutes(e *echo.Echo) {
	g := e.Group("/api/v1")

	g.POST("/login", loginApi)
	g.POST("/logout", logoutApi, authRequiredMiddleware)
	g.GET("/stats", serverStats, authRequiredMiddleware)
	g.GET("/tunnel-users", getTunnelUsers, authRequiredMiddleware)
	g.POST("/tunnel-users", createTunnelUser, authRequiredMiddleware)
	g.PUT("/tunnel-users", rotateTunnelUserSecretKey, authRequiredMiddleware)
}

func loginApi(c echo.Context) error {
	var p LoginPayload

	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	app := c.Get("app").(*app.App)

	token, err := app.User.Login(c.Request().Context(), p.Email, p.Password)
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

func logoutApi(c echo.Context) error {
	sessionCookie, err := c.Request().Cookie("beaver_session")

	if err != nil || sessionCookie == nil || sessionCookie.Value == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid session"})
	}

	app := c.Get("app").(*app.App)

	err = app.User.Logout(c.Request().Context(), sessionCookie.Value)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	sessionCookie.Value = ""
	sessionCookie.Expires = time.Now()
	sessionCookie.Path = "/"

	c.SetCookie(sessionCookie)

	return nil
}

func serverStats(c echo.Context) error {
	var result = make(map[string]any)

	v, _ := mem.VirtualMemory()
	result["memory_used"] = v.UsedPercent

	_c, _ := cpu.Percent(0, false)
	result["cpu_used"] = _c

	app := c.Get("app").(*app.App)
	app.Server.Lock.Lock()
	defer app.Server.Lock.Unlock()

	result["active_connections"] = len(app.Server.Pools)

	c.JSON(200, result)
	return nil
}

type createTunnelUserPayload struct {
	Email string
}

func createTunnelUser(c echo.Context) error {
	var payload createTunnelUserPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	app := c.Get("app").(*app.App)

	tunnelUser, err := app.User.CreateTunnelUser(c.Request().Context(), payload.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tunnelUser)
}

func getTunnelUsers(c echo.Context) error {
	app := c.Get("app").(*app.App)
	tunnelUsers, err := app.User.ListTunnelUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tunnelUsers)
}

func rotateTunnelUserSecretKey(c echo.Context) error {
	var payload createTunnelUserPayload
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid payload"})
	}

	app := c.Get("app").(*app.App)
	tunnelUsers, err := app.User.RotateTunnelUserSecretKey(c.Request().Context(), payload.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tunnelUsers)
}

func GetAdminHandler(app *app.App) *echo.Echo {
	adminRouter := echo.New()

	adminRouter.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	requiresLogin := func(c echo.Context, next echo.HandlerFunc) error {
		if c.Path() == "/dashboard" {
			return c.Redirect(307, "/")
		} else {
			return next(c)
		}
	}

	// Redirect non-subdomain pages based on valid session token
	adminRouter.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/" || c.Path() == "/dashboard" {
				err := authRequired(c)
				if err != nil {
					return requiresLogin(c, next)
				}
				if c.Path() == "/" {
					return c.Redirect(307, "/dashboard")
				}
			}
			return next(c)
		}
	})

	// Index
	if debug := os.Getenv("DEBUG"); debug == "True" {
		adminRouter.File("/", "./internal/server/templates/index.html")
		adminRouter.File("/dashboard", "./internal/server/templates/index.html")
	} else {
		fsysAssets, err := fs.Sub(web.DistAssets, "dist")
		if err != nil {
			panic(err)
		}
		AssetsHandler := http.FileServer(http.FS(fsysAssets))
		adminRouter.GET("/", func(c echo.Context) error {
			return c.Blob(http.StatusOK, "text/html", web.DistIndex)
		})
		adminRouter.GET("/dashboard", func(c echo.Context) error {
			return c.Blob(http.StatusOK, "text/html", web.DistIndex)
		})
		adminRouter.GET("/assets/*", echo.WrapHandler(AssetsHandler))
	}

	adminRouter.GET("/static/favicon.ico", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "image/x-icon", static.Favicon)
	})

	adminRouter.GET("/register", register)
	adminRouter.GET("/status", status)

	// Setup API routes
	setupApiRoutes(adminRouter)

	return adminRouter
}
