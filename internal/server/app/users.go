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

	if err := u.Store.FindOne(&superuser, badgerhold.Where("Email").Eq(email).Limit(1)); err != nil {
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
