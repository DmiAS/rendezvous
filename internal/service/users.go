package service

import (
	"context"
	"fmt"

	"rendezvous/internal/model"
)

type UserRepository interface {
	GetUsers(ctx context.Context) (model.Users, error)
	CreateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, login string) error
}
type UserService struct {
	r UserRepository
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUsers(ctx context.Context) (model.Users, error) {
	users, err := u.r.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failure to retrive users data from db: %s", err)
	}
	return users, nil
}

func (u *UserService) AddUser(ctx context.Context, user *model.User) error {
	if err := u.r.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failure to add user %+v: %s", user, err)
	}
	return nil
}

func (u *UserService) DeleteUser(ctx context.Context, login string) error {
	if err := u.r.DeleteUser(ctx, login); err != nil {
		return fmt.Errorf("failure to delete user %s: %s", login, err)
	}
	return nil
}
