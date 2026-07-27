package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

var requestValidator = validator.New()

type validatable interface {
	Validate(requestValidator *validator.Validate) error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrorInvalidArgument,
		)
	}

	var (
		err error
	)

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate(requestValidator)
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w",
			err,
			core_errors.ErrorInvalidArgument,
		)
	}

	return nil
}
