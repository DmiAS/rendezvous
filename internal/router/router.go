package router

import (
	"context"

	"rendezvous/internal/model"
)

type UserService interface {
	GetUsers(ctx context.Context) (model.Users, error)
}
type Router struct {
	u UserService
}

func NewRouter() *Router {
	return &Router{}
}
