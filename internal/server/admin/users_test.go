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

var store = newTestStore()

func TestCreateSuperUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	var err error

	ctx := context.Background()
	user := NewUserService(store)

	// No error while creating superuser
	_, err = user.CreateSuperUser(ctx, "test@beaver.com", "password")

	assert.NoError(t, err)

	// Creating superuser with duplicate email should throw error
	_, err = user.CreateSuperUser(ctx, "test@beaver.com", "password")
	assert.Error(t, err)
	assert.Equal(t, err, ErrDuplicateSuperUser)
}

func TestLoginSuperUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	_, _ = user.CreateSuperUser(ctx, "test@beaver.com", "password")

	token, _ := user.Login(ctx, "test@beaver.com", "password")

	superUser, _ := user.findUserByEmail(ctx, "test@beaver.com")
	assert.NotEqual(t, superUser.SessionToken, "")
	assert.Equal(t, superUser.SessionToken, token)

	token, err := user.Login(ctx, "test2@beaver.com", "password")
	assert.Error(t, err)
	assert.Equal(t, ErrWrongEmailOrPassword, err)
	assert.Equal(t, token, "")
}

func TestValidateSession(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	superUser, _ := user.CreateSuperUser(ctx, "test@beaver.com", "password")

	token, _ := user.Login(ctx, "test@beaver.com", "password")

	superUser2, err := user.ValidateSession(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, superUser.Email, superUser2.Email)

	s, err := user.ValidateSession(ctx, "random_token")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUserSession, err)
}

func TestLogoutSuperUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	_, _ = user.CreateSuperUser(ctx, "test@beaver.com", "password")

	token, _ := user.Login(ctx, "test@beaver.com", "password")

	err := user.Logout(ctx, token)
	assert.NoError(t, err)

	superUser, err := user.findUserByEmail(ctx, "test@beaver.com")
	assert.Equal(t, superUser.SessionToken, "")
	assert.NoError(t, err)

	err = user.Logout(ctx, "test2@beaver.com")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUserSession, err)
}
