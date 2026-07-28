package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/wasstend/todoapp-golang/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	APIVersion1 = APIVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(
	apiVersion APIVersion,
	middleware ...core_http_middleware.Middleware,
) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
		middleware: middleware,
	}
}

func (r *APIVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		r,
		r.middleware...,
	)
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		// pattern - "GET /tasks"
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())

	}
}
