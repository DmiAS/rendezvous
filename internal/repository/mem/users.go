package mem

import (
	"fmt"
	"sync"

	"github.com/DmiAS/rendezvous/internal/model"
)

type UsersRepository struct {
	m sync.Map
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{m: sync.Map{}}
}

func (u *UsersRepository) GetUsers() model.Users {
	// get all users who are not communicating right now
	var users model.Users
	u.m.Range(
		func(_, value interface{}) bool {
			user := value.(model.User)
			if !user.Blocked {
				users = append(users, user)
			}
			return true
		},
	)
	return users
}

func (u *UsersRepository) GetUser(login string) (*model.User, error) {
	v, ok := u.m.Load(login)
	if !ok {
		return nil, fmt.Errorf("no such user")
	}
	user := v.(model.User)
	return &user, nil
}

func (u *UsersRepository) UpdateUser(user *model.User) error {
	v, ok := u.m.Load(user.Name)
	if !ok {
		return fmt.Errorf("no such user")
	}

	currentUser := v.(model.User)
	if user.LocalAddress != "" {
		currentUser.LocalAddress = user.LocalAddress
	}
	if user.GlobalAddress != "" {
		currentUser.GlobalAddress = user.GlobalAddress
	}
	currentUser.Blocked = user.Blocked

	u.m.Store(user.Name, currentUser)
	return nil
}

func (u *UsersRepository) CreateUser(user *model.User) {
	u.m.Store(user.Name, *user)
}

func (u *UsersRepository) DeleteUser(login string) {
	u.m.Delete(login)
}
