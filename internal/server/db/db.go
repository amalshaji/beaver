package db

import (
	"log"
	"os"

	"github.com/amalshaji/beaver/internal/server/admin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewStore() *gorm.DB {
	// create database directory if not exists
	if _, err := os.Stat("./data"); err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("./data", os.ModePerm)
		}
	}
	db, err := gorm.Open(sqlite.Open("./data/beaver.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// should automigrate here?
	db.AutoMigrate(&admin.AdminUser{}, &admin.TunnelUser{}, &admin.Session{})

	return db
}
