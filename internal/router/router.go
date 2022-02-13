package router

import "github.com/DmiAS/rendezvous/internal/model"

type UserService interface {
	GetUsers() *model.InnerUsers
}
type Router struct {
	u UserService
}

func NewRouter(u UserService) *Router {
	return &Router{u: u}
}
