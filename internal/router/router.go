package router

import (
	"context"

	"github.com/DmiAS/rendezvous/internal/model"
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
