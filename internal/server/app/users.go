package app

import (
	"context"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	DB *bolt.DB
}

func NewUser(db *bolt.DB) *User {
	return &User{DB: db}
}

func (u *User) CreateSuperUser(ctx context.Context) {

}
