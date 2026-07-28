package core_http_middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(
	h http.Handler,
	m ...Middleware,
) http.Handler {
	if len(m) == 0 {
		return h
	}

	// for i := len(m) - 1; i >= 0; i-- {
	// 	h = m[i](h)
	// }

	for _, v := range slices.Backward(m) {
		h = v(h)
	}

	return h
}
