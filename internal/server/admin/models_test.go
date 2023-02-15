package admin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdminUserModel(t *testing.T) {
	adminuser := AdminUser{
		Email: "test@beaver.com",
	}

	adminuser.SetPassword("password")

	assert.NotEqual(t, "password", adminuser.PasswordHash)
	assert.NoError(t, adminuser.CheckPassword("password"))
	assert.Error(t, adminuser.CheckPassword("wrongpassword"))
}
