package db

import (
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

func NewDatabase() *bolt.DB {
	// Make sure data directory exists
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", os.ModePerm)
	}
	db, err := bolt.Open("./data/beaver.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
