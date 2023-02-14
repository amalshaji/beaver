package admin

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/timshannon/badgerhold/v4"
)

func newTestStore() *badgerhold.Store {
	options := badgerhold.DefaultOptions
	options.Dir = "./testdata"
	options.ValueDir = "./testdata"
	options.Logger = nil

	store, err := badgerhold.Open(options)
	if err != nil {
		log.Fatal(err)
	}

	return store
}

func TestCreateSuperUser(t *testing.T) {
	store := newTestStore()
	defer func() {
		store.Badger().DropAll()
	}()

	var err error

	ctx := context.Background()
	user := NewUserService(store)

	// No error while creating superuser
	err = user.CreateSuperUser(ctx, "test@beaver.com", "password")
	assert.NoError(t, err)

	// Creating superuser with duplicate email should throw error
	err = user.CreateSuperUser(ctx, "test@beaver.com", "password")
	assert.Equal(t, err.Error(), "superuser with the same email exists")
}
