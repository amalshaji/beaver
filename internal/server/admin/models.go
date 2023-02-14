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

type SuperUser struct {
	BaseModel

	Email        string `badgerhold:"unique"`
	PasswordHash string
	SessionToken string
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

func (s *SuperUser) CheckPassword(rawPassword string) error {
	rawPassword = utils.SanitizeString(rawPassword)
	err := bcrypt.CompareHashAndPassword([]byte(s.PasswordHash), []byte(rawPassword))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}

func (s *SuperUser) GenerateSessionToken() error {
	s.SessionToken = utils.GenerateUUIDV4().String()
	return nil
}

func (s *SuperUser) ResetSessionToken() error {
	s.SessionToken = ""
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
