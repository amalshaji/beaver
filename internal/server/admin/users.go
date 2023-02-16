package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/amalshaji/beaver/internal/utils"
	"github.com/timshannon/badgerhold/v4"
)

var ErrAdminUserNotFound = errors.New("admin user does not exist")
var ErrTunnelUserNotFound = errors.New("tunnel user does not exist")
var ErrInvalidUserSession = errors.New("invalid user session")
var ErrWrongEmailOrPassword = errors.New("wrong email or password")
var ErrDuplicateAdminUser = errors.New("admin user with the same email exists")
var ErrDuplicateTunnelUser = errors.New("tunnel user with the same email exists")

type User struct {
	Store *badgerhold.Store
}

func NewUserService(store *badgerhold.Store) *User {
	return &User{Store: store}
}

func (u *User) findUserByEmail(ctx context.Context, email string) (*AdminUser, error) {
	email = utils.SanitizeString(email)

	var superUser AdminUser
	if err := u.Store.FindOne(&superUser, badgerhold.Where("Email").Eq(email)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, ErrAdminUserNotFound
		}
		return nil, err
	}
	return &superUser, nil
}

func (u *User) CreateUser(ctx context.Context, email, password string, isSuperUser bool) (*AdminUser, error) {
	email = utils.SanitizeString(email)
	password = utils.SanitizeString(password)

	existingAdminUser, err := u.findUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrAdminUserNotFound) {
		return nil, err
	}

	if existingAdminUser != nil {
		return nil, ErrDuplicateAdminUser
	}

	var adminUser AdminUser

	adminUser.Email = email
	adminUser.SetPassword(password)
	adminUser.IsSuperUser = isSuperUser
	adminUser.MarkAsNew()

	if err := u.Store.Insert(badgerhold.NextSequence(), adminUser); err != nil {
		if errors.Is(err, badgerhold.ErrUniqueExists) {
			return nil, ErrDuplicateAdminUser
		}
		return nil, err
	}

	return &adminUser, nil
}

func (u *User) CreateAdminUser(ctx context.Context, email, password string) (*AdminUser, error) {
	return u.CreateUser(ctx, email, password, false)
}

func (u *User) CreateSuperUser(ctx context.Context, email, password string) (*AdminUser, error) {
	return u.CreateUser(ctx, email, password, true)
}

func (u *User) Login(ctx context.Context, email, password string) (string, error) {
	email = utils.SanitizeString(email)
	password = utils.SanitizeString(password)

	var adminUser *AdminUser

	adminUser, err := u.findUserByEmail(ctx, email)
	if err != nil && errors.Is(err, ErrAdminUserNotFound) {
		return "", ErrWrongEmailOrPassword
	}

	if err := adminUser.CheckPassword(password); err != nil {
		return "", ErrWrongEmailOrPassword
	}

	adminUser.GenerateSessionToken()

	u.Store.UpdateMatching(&AdminUser{}, badgerhold.Where("Email").Eq(email), func(record interface{}) error {
		update, ok := record.(*AdminUser)
		if !ok {
			return fmt.Errorf("error while updating superuser")
		}
		update.SessionToken = adminUser.SessionToken
		return nil
	})

	return adminUser.SessionToken, nil
}

func (u *User) Logout(ctx context.Context, sessionToken string) error {
	var err error

	if _, err = u.ValidateSession(ctx, sessionToken); err != nil {
		return err
	}

	u.Store.UpdateMatching(&AdminUser{}, badgerhold.Where("SessionToken").Eq(sessionToken), func(record interface{}) error {
		update, ok := record.(*AdminUser)
		if !ok {
			return fmt.Errorf("error while updating superuser")
		}

		update.SessionToken = ""
		return nil
	})

	return nil
}

func (u *User) ValidateSession(ctx context.Context, sessionToken string) (*AdminUser, error) {
	var adminUser AdminUser

	if err := u.Store.FindOne(&adminUser, badgerhold.Where("SessionToken").Eq(sessionToken)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, ErrInvalidUserSession
		}
		return nil, err
	}
	return &adminUser, nil
}

func (u *User) findTunnelUserByEmail(ctx context.Context, email string) (*TunnelUser, error) {
	email = utils.SanitizeString(email)

	var tunnelUser TunnelUser

	if err := u.Store.FindOne(&tunnelUser, badgerhold.Where("Email").Eq(email)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, ErrTunnelUserNotFound
		}
		return nil, err
	}
	return &tunnelUser, nil
}

func (u *User) CreateTunnelUser(ctx context.Context, email string) (*TunnelUser, error) {
	email = utils.SanitizeString(email)

	existingTunnelUser, err := u.findTunnelUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrTunnelUserNotFound) {
		return nil, err
	}

	if existingTunnelUser != nil {
		return nil, ErrDuplicateTunnelUser
	}

	if err := utils.ValidateEmail(email); err != nil {
		return nil, fmt.Errorf("enter a valid email address")
	}

	var tunnelUser TunnelUser

	tunnelUser.Email = email
	tunnelUser.RotateSecretKey()
	tunnelUser.MarkAsNew()

	if err := u.Store.Insert(badgerhold.NextSequence(), tunnelUser); err != nil {
		if errors.Is(err, badgerhold.ErrUniqueExists) {
			return nil, ErrDuplicateTunnelUser
		}
		return nil, err
	}

	return &tunnelUser, nil
}

func (u *User) GetTunnelUserBySecret(ctx context.Context, secretKey string) (*TunnelUser, error) {
	secretKey = utils.SanitizeString(secretKey)

	var tunnelUser TunnelUser

	if err := u.Store.FindOne(&tunnelUser, badgerhold.Where("SecretKey").Eq(secretKey)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, ErrTunnelUserNotFound
		}
		return nil, err
	}
	return &tunnelUser, nil
}

func (u *User) ListTunnelUsers(ctx context.Context) ([]TunnelUser, error) {
	var tunnelUsers []TunnelUser

	if err := u.Store.Find(&tunnelUsers, nil); err != nil {
		return nil, err
	}
	if tunnelUsers == nil {
		return []TunnelUser{}, nil
	}
	return tunnelUsers, nil
}
