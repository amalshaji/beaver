package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/timshannon/badgerhold/v4"
)

var ErrSuperUserNotFound = errors.New("user does not exist")
var ErrInvalidUserSession = errors.New("invalid user session")
var ErrWrongEmailOrPassword = errors.New("wrong email or password")
var ErrDuplicateSuperUser = errors.New("superuser with the same email exists")

type User struct {
	Store *badgerhold.Store
}

func NewUserService(store *badgerhold.Store) *User {
	return &User{Store: store}
}

func (u *User) findUserByEmail(ctx context.Context, email string) (*SuperUser, error) {
	var superUser SuperUser
	if err := u.Store.FindOne(&superUser, badgerhold.Where("Email").Eq(email)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, ErrSuperUserNotFound
		}
		return nil, err
	}
	return &superUser, nil
}

func (u *User) CreateSuperUser(ctx context.Context, email, password string) (*SuperUser, error) {
	existingSuperUser, err := u.findUserByEmail(ctx, email)
	fmt.Printf("existing user: %#v\n", existingSuperUser)

	if err != nil && !errors.Is(err, ErrSuperUserNotFound) {
		return nil, err
	}

	if existingSuperUser != nil {
		return nil, ErrDuplicateSuperUser
	}

	var superUser SuperUser

	superUser.Email = email
	superUser.SetPassword(password)
	superUser.MarkAsNew()

	if err := u.Store.Insert(badgerhold.NextSequence(), superUser); err != nil {
		if errors.Is(err, badgerhold.ErrUniqueExists) {
			return nil, ErrDuplicateSuperUser
		}
		return nil, err
	}

	return &superUser, nil
}

func (u *User) Login(ctx context.Context, email, password string) (string, error) {
	var superUser *SuperUser

	superUser, err := u.findUserByEmail(ctx, email)
	if err != nil && errors.Is(err, ErrSuperUserNotFound) {
		return "", ErrWrongEmailOrPassword
	}

	if err := superUser.CheckPassword(password); err != nil {
		return "", ErrWrongEmailOrPassword
	}

	superUser.GenerateSessionToken()

	u.Store.UpdateMatching(&SuperUser{}, badgerhold.Where("Email").Eq(email), func(record interface{}) error {
		update, ok := record.(*SuperUser)
		if !ok {
			return fmt.Errorf("error while updating superuser")
		}
		update.SessionToken = superUser.SessionToken
		return nil
	})

	return superUser.SessionToken, nil
}

func (u *User) Logout(ctx context.Context, sessionToken string) error {
	var err error

	if _, err = u.ValidateSession(ctx, sessionToken); err != nil {
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

func (u *User) ValidateSession(ctx context.Context, sessionToken string) (*SuperUser, error) {
	var superuser SuperUser

	if err := u.Store.FindOne(&superuser, badgerhold.Where("SessionToken").Eq(sessionToken)); err != nil {
		if errors.Is(err, badgerhold.ErrNotFound) {
			return nil, ErrInvalidUserSession
		}
		return nil, err
	}
	return &superuser, nil
}
