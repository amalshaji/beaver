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
	superUser, err := user.CreateSuperUser(ctx, "test@beaver.com", "password")
	assert.NoError(t, err)
	assert.True(t, superUser.IsSuperUser)

	// Creating multiple superusers should fail
	_, err = user.CreateSuperUser(ctx, "test@beaver.com", "password")
	assert.Error(t, err)
	assert.Equal(t, err, ErrMultipleSuperuserError)
}

func TestAdminSuperUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	var err error

	ctx := context.Background()
	user := NewUserService(store)

	// No error while creating adminuser
	adminUser, err := user.CreateAdminUser(ctx, "test@beaver.com", "password")
	assert.NoError(t, err)
	assert.False(t, adminUser.IsSuperUser)

	// Creating adminuser with duplicate email should throw error
	_, err = user.CreateAdminUser(ctx, "test@beaver.com", "password")
	assert.Error(t, err)
	assert.Equal(t, err, ErrDuplicateAdminUser)

	// Creating superuser with duplicate email should throw error
	_, err = user.CreateSuperUser(ctx, "test@beaver.com", "password")
	assert.Error(t, err)
	assert.Equal(t, err, ErrDuplicateAdminUser)
}

func TestLoginAdminUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	_, _ = user.CreateAdminUser(ctx, "test@beaver.com", "password")

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

	superUser, _ := user.CreateAdminUser(ctx, "test@beaver.com", "password")

	token, _ := user.Login(ctx, "test@beaver.com", "password")

	superUser2, err := user.ValidateSession(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, superUser.Email, superUser2.Email)

	s, err := user.ValidateSession(ctx, "random_token")
	assert.Nil(t, s)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidUserSession, err)
}

func TestLogoutAdminUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	_, _ = user.CreateAdminUser(ctx, "test@beaver.com", "password")

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

func TestCreateTunnelUser(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	tu, err := user.CreateTunnelUser(ctx, "test@beaver.com")
	assert.NoError(t, err)
	assert.Equal(t, "test@beaver.com", tu.Email)
	assert.NotEqual(t, "", tu.SecretKey)
}

func TestGetTunnelUserBySecretKey(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	tu, _ := user.CreateTunnelUser(ctx, "test@beaver.com")

	ntu, err := user.GetTunnelUserBySecret(ctx, tu.SecretKey)
	assert.NoError(t, err)
	assert.Equal(t, tu.Email, ntu.Email)
}

func TestRotateTunnelUserSecretKey(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	tu, _ := user.CreateTunnelUser(ctx, "test@beaver.com")

	_, err := user.RotateTunnelUserSecretKey(ctx, tu.Email)
	assert.NoError(t, err)

	ntu, _ := user.findTunnelUserByEmail(ctx, "test@beaver.com")
	assert.NotEqual(t, tu.SecretKey, ntu.SecretKey)

	nontu, err := user.RotateTunnelUserSecretKey(ctx, "test2@beaver.com")
	assert.Error(t, err)
	assert.Equal(t, ErrTunnelUserNotFound, err)
	assert.Nil(t, nontu)
}

func TestListTunnelUsers(t *testing.T) {
	defer func() {
		store.Badger().DropAll()
	}()

	ctx := context.Background()
	user := NewUserService(store)

	tunnelUsers, err := user.ListTunnelUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(tunnelUsers))

	_, _ = user.CreateTunnelUser(ctx, "test@beaver.com")
	_, _ = user.CreateTunnelUser(ctx, "test2@beaver.com")
	_, _ = user.CreateTunnelUser(ctx, "test3@beaver.com")

	tunnelUsers, err = user.ListTunnelUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(tunnelUsers))
	assert.Equal(t, "test@beaver.com", tunnelUsers[0].Email)
	assert.Equal(t, "test2@beaver.com", tunnelUsers[1].Email)
	assert.Equal(t, "test3@beaver.com", tunnelUsers[2].Email)
}
