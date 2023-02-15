package admin

import (
	"errors"
	"time"

	"github.com/amalshaji/beaver/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var ErrWrongPassword = errors.New("wrong password")

type BaseModel struct {
	ID        uint64 `badgerhold:"key"`
	CreatedAt time.Time
}

func (b *BaseModel) MarkAsNew() {
	b.CreatedAt = time.Now()
}

type AdminUser struct {
	BaseModel

	Email        string `badgerhold:"unique"`
	PasswordHash string
	SessionToken string
	IsSuperUser  bool
}

func (a *AdminUser) SetPassword(rawPassword string) error {
	rawPassword = utils.SanitizeString(rawPassword)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), 13)
	if err != nil {
		return err
	}
	a.PasswordHash = string(hashedPassword)
	return nil
}

func (a *AdminUser) CheckPassword(rawPassword string) error {
	rawPassword = utils.SanitizeString(rawPassword)
	err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte(rawPassword))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}

func (a *AdminUser) GenerateSessionToken() error {
	a.SessionToken = utils.GenerateUUIDV4().String()
	return nil
}

func (a *AdminUser) ResetSessionToken() error {
	a.SessionToken = ""
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
