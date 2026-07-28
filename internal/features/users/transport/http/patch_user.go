package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_types "github.com/wasstend/todoapp-golang/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate(requestValidator *validator.Validate) error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("FullName length must be minimum 3 and maximum 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			if err := requestValidator.Var(*r.PhoneNumber.Value, "e164"); err != nil {
				return fmt.Errorf("PhoneNumber must be in format e164")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	logger.Debug("invoke PatchUser Handler")

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id path value",
		)
		return
	}

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate http request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusOK)

}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
