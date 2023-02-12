package app

import "github.com/timshannon/badgerhold/v4"

type Dashboard struct {
	Store *badgerhold.Store
}

func NewDashboardService(store *badgerhold.Store) *Dashboard {
	return &Dashboard{Store: store}
}
