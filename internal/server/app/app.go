package app

import (
	"github.com/amalshaji/beaver/internal/server/admin"
	"github.com/amalshaji/beaver/internal/server/db"
	"github.com/amalshaji/beaver/internal/server/tunnel"
	"github.com/timshannon/badgerhold/v4"
)

type App struct {
	Store     *badgerhold.Store
	Dashboard *admin.Dashboard
	User      *admin.User
	Server    *tunnel.Server
}

func NewApp(configFile string) *App {
	store := db.NewStore()
	return &App{
		Store:     store,
		Dashboard: admin.NewDashboardService(store),
		User:      admin.NewUserService(store),
		Server:    tunnel.NewServer(configFile),
	}
}

func (app *App) Start() {
	app.Server.Start()
}

func (app *App) Shutdown() {
	// Shutdown the tunnel server
	app.Server.Shutdown()

	// Close database connection
	app.Store.Close()
}
