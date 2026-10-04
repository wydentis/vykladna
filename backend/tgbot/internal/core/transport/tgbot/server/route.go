package core_tgbot_server

import (
	"context"
	"reflect"

	core_tgbot_middleware "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/middleware"
	core_tgbot_types "github.com/wydentis/vykladna/tgbot/internal/core/transport/tgbot/types"
)

type Route struct {
	payload     reflect.Type
	Handler     core_tgbot_middleware.HandlerFunc
	Middlewares []core_tgbot_middleware.Middleware
}

func NewRoute[T any](h func(context.Context, T) error, m ...core_tgbot_middleware.Middleware) Route {
	return Route{
		payload: reflect.TypeFor[T](),
		Handler: func(ctx context.Context, u core_tgbot_types.Update) error {
			return h(ctx, u.Payload.(T))
		},
		Middlewares: m,
	}
}

func (r *Route) WithMiddlewares() core_tgbot_middleware.HandlerFunc {
	return core_tgbot_middleware.Chain(r.Handler, r.Middlewares...)
}
