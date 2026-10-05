package core_tgbot_server

import (
	"reflect"

	core_tgbot_middleware "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/middleware"
)

type Router struct {
	handlers    map[reflect.Type]core_tgbot_middleware.HandlerFunc
	middlewares []core_tgbot_middleware.Middleware
}

func NewRouter(middlewares ...core_tgbot_middleware.Middleware) *Router {
	return &Router{
		handlers:    map[reflect.Type]core_tgbot_middleware.HandlerFunc{},
		middlewares: middlewares,
	}
}

func (r *Router) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		if _, dup := r.handlers[route.payload]; dup {
			panic("tgbot router: duplicate route for " + route.payload.String())
		}
		r.handlers[route.payload] = route.WithMiddlewares()
	}
}
