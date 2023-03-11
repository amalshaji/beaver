package admin

import (
	"context"
	"errors"
	"fmt"

	"github.com/amalshaji/beaver/internal/utils"
	"gorm.io/gorm"
)

var ErrAdminUserNotFound = errors.New("admin user does not exist")
var ErrTunnelUserNotFound = errors.New("tunnel user does not exist")
var ErrInvalidUserSession = errors.New("invalid user session")
var ErrWrongEmailOrPassword = errors.New("wrong email or password")
var ErrDuplicateAdminUser = errors.New("admin user with the same email exists")
var ErrDuplicateTunnelUser = errors.New("tunnel user with the same email exists")
var ErrMultipleSuperuserError = errors.New("you cannot create more than one superuser")

type UserService struct {
	DB *gorm.DB
}

func NewUserService(store *gorm.DB) *UserService {
	return &UserService{DB: store}
}

func (u *UserService) findUserByEmail(ctx context.Context, email string) (*AdminUser, error) {
	email = utils.SanitizeString(email)

	var user AdminUser
	result := u.DB.Where(&AdminUser{Email: email}).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrAdminUserNotFound
		}
		return nil, result.Error
	}

	return &user, nil
}

func (u *UserService) CreateUser(ctx context.Context, email, password string, superUser bool) (*AdminUser, error) {
	email = utils.SanitizeString(email)
	password = utils.SanitizeString(password)

	existingAdminUser, err := u.findUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrAdminUserNotFound) {
		return nil, err
	}

	if existingAdminUser != nil {
		return nil, ErrDuplicateAdminUser
	}

	adminUser := AdminUser{
		Email:     email,
		SuperUser: superUser,
	}
	adminUser.SetPassword(password)

	result := u.DB.Create(&adminUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &adminUser, nil
}

func (u *UserService) CreateAdminUser(ctx context.Context, email, password string) (*AdminUser, error) {
	return u.CreateUser(ctx, email, password, false)
}

func (u *UserService) CanCreateSuperUser(ctx context.Context) error {
	var count int64

	result := u.DB.Model(&AdminUser{}).Where("super_user = ?", true).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	if count == 0 {
		return nil
	}

	return ErrMultipleSuperuserError
}

func (u *UserService) CreateSuperUser(ctx context.Context, email, password string) (*AdminUser, error) {
	if err := u.CanCreateSuperUser(ctx); err != nil {
		return nil, err
	}
	return u.CreateUser(ctx, email, password, true)
}

func (u *UserService) Login(ctx context.Context, email, password string) (string, error) {
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

	session := Session{
		Token: utils.GenerateSessionToken(),
	}
	adminUser.Session = session

	result := u.DB.Save(&adminUser)
	if result.Error != nil {
		return "", result.Error
	}

	return adminUser.Session.Token, nil
}

func (u *UserService) Logout(ctx context.Context, sessionToken string) error {
	result := u.DB.Where(&Session{Token: sessionToken}).Delete(&Session{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrInvalidUserSession
	}

	return nil
}

func (u *UserService) ValidateSession(ctx context.Context, sessionToken string) (*AdminUser, error) {
	var adminUser AdminUser

	result := u.DB.
		Joins("JOIN sessions on sessions.admin_user_id = admin_users.id").
		Where("sessions.token = ?", sessionToken).First(&adminUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidUserSession
		}
		return nil, result.Error
	}

	return &adminUser, nil
}

func (u *UserService) findTunnelUserByEmail(ctx context.Context, email string) (*TunnelUser, error) {
	email = utils.SanitizeString(email)

	var tunnelUser TunnelUser

	result := u.DB.Where(&TunnelUser{Email: email}).First(&tunnelUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTunnelUserNotFound
		}
		return nil, result.Error
	}

	return &tunnelUser, nil
}

func (u *UserService) CreateTunnelUser(ctx context.Context, email string) (*TunnelUser, error) {
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

	tunnelUser := TunnelUser{
		Email: email,
	}
	tunnelUser.RotateSecretKey()

	result := u.DB.Save(&tunnelUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tunnelUser, nil
}

func (u *UserService) GetTunnelUserBySecret(ctx context.Context, secretKey string) (*TunnelUser, error) {
	secretKey = utils.SanitizeString(secretKey)

	var tunnelUser TunnelUser

	result := u.DB.Where(&TunnelUser{SecretKey: &secretKey}).First(&tunnelUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrTunnelUserNotFound
		}
		return nil, result.Error
	}

	return &tunnelUser, nil
}

func (u *UserService) ListTunnelUsers(ctx context.Context) ([]TunnelUser, error) {
	var tunnelUsers []TunnelUser

	result := u.DB.Find(&tunnelUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	if len(tunnelUsers) == 0 {
		return []TunnelUser{}, nil
	}

	return tunnelUsers, nil
}

func (u *UserService) RotateTunnelUserSecretKey(ctx context.Context, email string) (*TunnelUser, error) {
	tunnelUser, err := u.findTunnelUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	tunnelUser.RotateSecretKey()

	result := u.DB.Save(&tunnelUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return tunnelUser, nil
}
