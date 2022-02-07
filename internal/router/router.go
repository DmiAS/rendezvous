package router

import (
	"context"

	"rendezvous/internal/model"
)

type UserService interface {
	GetUsers(ctx context.Context) (model.InnerUsers, error)
}
type Router struct {
	u UserService
}

func NewRouter(u UserService) *Router {
	return &Router{u: u}
}
