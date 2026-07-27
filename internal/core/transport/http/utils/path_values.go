package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

func GetIntPathValue(r *http.Request, key string) (int, error) {
	path_value := r.PathValue(key)
	if path_value == "" {
		return 0, fmt.Errorf(
			"no key='%s' in path values: %w",
			key,
			core_errors.ErrorInvalidArgument,
		)
	}

	value, err := strconv.Atoi(path_value)
	if err != nil {
		return 0, fmt.Errorf(
			"path value='%s' by key='%s' not a valid integer: %v: %w",
			path_value,
			key,
			err,
			core_errors.ErrorInvalidArgument,
		)
	}

	return value, nil
}
