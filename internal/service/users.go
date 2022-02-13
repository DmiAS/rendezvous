package service

import (
	"context"
	"fmt"

	"github.com/DmiAS/rendezvous/internal/model"
)

type UserRepository interface {
	GetUsers(ctx context.Context) (model.Users, error)
	GetUser(ctx context.Context, login string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, login string) error
}
type UserService struct {
	r UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{r: r}
}

func (u *UserService) GetUsers(ctx context.Context) (*model.InnerUsers, error) {
	users, err := u.r.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failure to retrive users data from db: %s", err)
	}
	names := make([]string, 0, len(users))
	for i := range users {
		names = append(names, users[i].Name)
	}
	return &model.InnerUsers{Names: names}, nil
}

func (u *UserService) GetUser(ctx context.Context, login string) (*model.User, error) {
	user, err := u.r.GetUser(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failure to get uses %s: %s", login, err)
	}
	return user, nil
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

func (u *UserService) BlockUser(ctx context.Context, user string) error {
	if err := u.r.UpdateUser(ctx, &model.User{Name: user, Blocked: true}); err != nil {
		return fmt.Errorf("failure to update user with data %+v: %s", user, err)
	}
	return nil
}
