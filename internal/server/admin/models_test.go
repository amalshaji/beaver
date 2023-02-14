package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuperUserModel(t *testing.T) {
	superuser := SuperUser{
		Email: "test@beaver.com",
	}

	superuser.SetPassword("password")

	assert.NotEqual(t, "password", superuser.PasswordHash)
	assert.NoError(t, superuser.CheckPassword("password"))
	assert.Error(t, superuser.CheckPassword("wrongpassword"))
}
