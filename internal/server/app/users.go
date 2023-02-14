package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/timshannon/badgerhold/v4"
)

type User struct {
	Store *badgerhold.Store
}

func NewUserService(store *badgerhold.Store) *User {
	return &User{Store: store}
}

func (u *User) CreateSuperUser(ctx context.Context, email, password string) error {
	var superuser SuperUser

	if err := u.Store.FindOne(&superuser, badgerhold.Where("Email").Eq(email)); err != nil {
		if !errors.Is(err, badgerhold.ErrNotFound) {
			return err
		}
	}

	if superuser.Email != "" {
		return fmt.Errorf("superuser with the same email exists")
	}

	superuser.Email = email
	superuser.SetPassword(password)
	superuser.MarkAsNew()

	if err := u.Store.Insert(badgerhold.NextSequence(), superuser); err != nil {
		if errors.Is(err, badgerhold.ErrUniqueExists) {
			return fmt.Errorf("superuser with the same email exists")
		}
		return err
	}

	return nil
}

func (u *User) Login(email, password string) (string, error) {
	var superuser SuperUser

	if err := u.Store.FindOne(&superuser, badgerhold.Where("Email").Eq(email)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return "", fmt.Errorf("wrong email or password")
		}
		return "", err
	}

	if err := superuser.CheckPassword(password); err != nil {
		return "", fmt.Errorf("wrong email or password")
	}

	superuser.GenerateSessionToken()

	u.Store.UpdateMatching(&SuperUser{}, badgerhold.Where("Email").Eq(email), func(record interface{}) error {
		update, ok := record.(*SuperUser)
		if !ok {
			return fmt.Errorf("error while updating superuser")
		}
		update.SessionToken = superuser.SessionToken
		return nil
	})

	return superuser.SessionToken, nil
}

func (u *User) Logout(sessionToken string) error {
	var err error

	if _, err = u.ValidateSession(sessionToken); err != nil {
		return err
	}

	u.Store.UpdateMatching(&SuperUser{}, badgerhold.Where("SessionToken").Eq(sessionToken), func(record interface{}) error {
		update, ok := record.(*SuperUser)
		if !ok {
			return fmt.Errorf("error while updating superuser")
		}

		update.SessionToken = ""
		return nil
	})

	return nil
}

func (u *User) ValidateSession(sessionToken string) (*SuperUser, error) {
	var superuser SuperUser

	if err := u.Store.FindOne(&superuser, badgerhold.Where("SessionToken").Eq(sessionToken)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, fmt.Errorf("invalid user session")
		}
		return nil, err
	}
	return &superuser, nil
}
