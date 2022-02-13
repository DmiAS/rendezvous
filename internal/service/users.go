package service

import (
	"fmt"

	"github.com/DmiAS/rendezvous/internal/model"
)

type UserRepository interface {
	GetUsers() model.Users
	GetUser(login string) (*model.User, error)
	CreateUser(user *model.User)
	UpdateUser(user *model.User) error
	DeleteUser(login string)
}
type UserService struct {
	r UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{r: r}
}

func (u *UserService) GetUsers() *model.InnerUsers {
	users := u.r.GetUsers()
	names := make([]string, 0, len(users))
	for i := range users {
		names = append(names, users[i].Name)
	}
	return &model.InnerUsers{Names: names}
}

func (u *UserService) GetUser(login string) (*model.User, error) {
	return u.r.GetUser(login)
}

func (u *UserService) AddUser(user *model.User) {
	u.r.CreateUser(user)
}

func (u *UserService) DeleteUser(login string) {
	u.r.DeleteUser(login)
}

func (u *UserService) BlockUser(user string) error {
	if err := u.r.UpdateUser(&model.User{Name: user, Blocked: true}); err != nil {
		return fmt.Errorf("failure to update user with data %+v: %s", user, err)
	}
	return nil
}
