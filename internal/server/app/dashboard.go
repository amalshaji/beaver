package app

import bolt "go.etcd.io/bbolt"

type Dashboard struct {
	DB *bolt.DB
}

func NewDashboard(db *bolt.DB) *Dashboard {
	return &Dashboard{DB: db}
}
