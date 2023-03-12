package app

import (
	"github.com/amalshaji/beaver/internal/server/admin"
	"github.com/amalshaji/beaver/internal/server/db"
	"github.com/amalshaji/beaver/internal/server/tunnel"
	"gorm.io/gorm"
)

type App struct {
	DB     *gorm.DB
	User   *admin.UserService
	Server *tunnel.Server
}

func NewApp(configFile string) *App {
	db := db.NewStore()
	return &App{
		DB:     db,
		User:   admin.NewUserService(db),
		Server: tunnel.NewServer(configFile, db),
	}
}

func (app *App) Start() {
	app.Server.Start()
}

func (app *App) Shutdown() {
	// Shutdown the tunnel server
	app.Server.Shutdown()
}
