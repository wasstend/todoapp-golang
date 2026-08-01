package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

// getQueryParam reads the query parameter by key and parses it with parse.
// Returns nil if the parameter is absent, or a wrapped ErrorInvalidArgument
// if parsing fails.
func getQueryParam[T any](r *http.Request, key string, parse func(string) (T, error)) (*T, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	value, err := parse(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' is not valid: %v: %w",
			param,
			key,
			err,
			core_errors.ErrorInvalidArgument,
		)
	}

	return &value, nil
}

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	return getQueryParam(r, key, strconv.Atoi)
}

func GetTimeQueryParam(r *http.Request, key string, layout string) (*time.Time, error) {
	return getQueryParam(r, key, func(param string) (time.Time, error) {
		return time.Parse(layout, param)
	})
}
