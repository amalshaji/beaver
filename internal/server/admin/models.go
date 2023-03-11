package admin

import (
	"errors"
	"time"

	"github.com/amalshaji/beaver/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrWrongPassword  = errors.New("wrong password")
	ErrWrongSecretKey = errors.New("wrong secret key")
)

type AdminUser struct {
	gorm.Model

	Email        string `gorm:"index,unique"`
	PasswordHash string `gorm:"unique" json:"-"`
	SuperUser    bool

	Session Session
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

type Session struct {
	gorm.Model

	Token       string `gorm:"index, unique"`
	AdminUserId uint
}

func (s *Session) GenerateSessionToken() error {
	s.Token = utils.GenerateSessionToken()
	return nil
}

type TunnelUser struct {
	gorm.Model

	Email        string  `gorm:"index,unique"`
	SecretKey    *string `gorm:"index,unique" json:"-"`
	LastActiveAt *time.Time
}

func (t *TunnelUser) RotateSecretKey() string {
	newSecretKey := utils.GenerateSecretKey()
	newSecretKeyHashedBytes, _ := bcrypt.GenerateFromPassword([]byte(newSecretKey), 13)
	newSecretKeyHashed := string(newSecretKeyHashedBytes)
	t.SecretKey = &newSecretKeyHashed
	return newSecretKey
}

func (t *TunnelUser) ValidateSecretKey(secretKey string) error {
	secretKey = utils.SanitizeString(secretKey)
	err := bcrypt.CompareHashAndPassword([]byte(*t.SecretKey), []byte(secretKey))
	if err != nil {
		return ErrWrongPassword
	}
	return nil
}
