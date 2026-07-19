package router

import (
	"github.com/go-telegram/ui/keyboard/inline"
	"go_telegram_bot/internal/domain/entity"
)

type Router struct {
	routes map[entity.UserAction]inline.OnSelect
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[entity.UserAction]inline.OnSelect),
	}
}

func (r *Router) Register(action entity.UserAction, handler inline.OnSelect) {
	r.routes[action] = handler
}

func (r *Router) Route(action entity.UserAction) inline.OnSelect {
	h, ok := r.routes[action]
	if !ok {
		return nil
	}
	return h
}
