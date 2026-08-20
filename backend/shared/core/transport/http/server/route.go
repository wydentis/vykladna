package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
)

type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []core_http_middleware.Middleware
}

func (r *Route) WithMiddlewares() http.Handler {
	return core_http_middleware.Chain(r.Handler, r.Middlewares...)
}
