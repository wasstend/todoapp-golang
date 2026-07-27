package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrorInvalidArgument,
		)
	}

	return &value, nil
}
