package app

import (
	"time"

	"github.com/amalshaji/beaver/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type BaseModel struct {
	ID        uint64 `badgerhold:"key"`
	CreatedAt time.Time
}

func (b *BaseModel) MarkAsNew() {
	b.CreatedAt = time.Now()
}

type SuperUser struct {
	BaseModel

	Email        string `badgerhold:"unique"`
	PasswordHash string
}

func (s *SuperUser) SetPassword(rawPassword string) error {
	rawPassword = utils.SanitizeString(rawPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 13)
	if err != nil {
		return err
	}
	s.PasswordHash = string(hashedPassword)
	return nil
}

type UserSession struct {
	BaseModel

	Token string `badgerhold:"unique"`
	User  SuperUser
}

func (u *UserSession) GenerateSessionToken() error {
	u.Token = utils.GenerateUUIDV4().String()
	return nil
}

type TunnelUser struct {
	BaseModel

	Email     string `badgerhold:"unique"`
	SecretKey string
}

func (t *TunnelUser) RotateSecretKey() error {
	t.SecretKey = utils.GenerateUUIDV4().String()
	return nil
}
