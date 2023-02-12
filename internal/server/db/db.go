package db

import (
	"log"

	"github.com/timshannon/badgerhold/v4"
)

func NewStore() *badgerhold.Store {
	options := badgerhold.DefaultOptions
	options.Dir = "data"
	options.ValueDir = "data"
	options.Logger = nil

	store, err := badgerhold.Open(options)
	if err != nil {
		log.Fatal(err)
	}

	return store
}
