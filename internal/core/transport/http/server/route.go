package core_http_server

import (
	"net/http"

	core_http_middleware "github.com/wasstend/todoapp-golang/internal/core/transport/http/middleware"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []core_http_middleware.Middleware
}

func (route *Route) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(
		route.Handler,
		route.Middleware...,
	)
}
