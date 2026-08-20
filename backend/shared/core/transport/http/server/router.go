package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/wydentis/vykladna/shared/core/transport_http/middleware"
)

type APIVersion string

var (
	APIVersion1 APIVersion = "v1"
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion  APIVersion
	middlewares []core_http_middleware.Middleware
}

func NewAPIVersionRouter(apiVersion APIVersion, middlewares ...core_http_middleware.Middleware) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:    http.NewServeMux(),
		apiVersion:  apiVersion,
		middlewares: middlewares,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, route.WithMiddlewares())
	}
}

func (r *APIVersionRouter) WithMiddlewares() http.Handler {
	return core_http_middleware.Chain(r, r.middlewares...)
}
