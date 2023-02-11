package app

import (
	"github.com/amalshaji/beaver/internal/server/db"
	"github.com/amalshaji/beaver/internal/server/tunnel"

	bolt "go.etcd.io/bbolt"
)

type App struct {
	DB        *bolt.DB
	Dashboard *Dashboard
	User      *User
	Server    *tunnel.Server
}

func NewApp(configFile string) *App {
	db := db.NewDatabase()
	return &App{
		DB:        db,
		Dashboard: NewDashboard(db),
		User:      NewUser(db),
		Server:    tunnel.NewServer(configFile),
	}
}

func (app *App) Shutdown() {
	// Shutdown the tunnel server
	app.Server.Shutdown()

	// Close database connection
	app.DB.Close()
}
