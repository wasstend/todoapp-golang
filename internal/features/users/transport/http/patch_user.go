package users_transport_http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_types "github.com/wasstend/todoapp-golang/internal/core/transport/http/types"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"John Doe"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79991234567"`
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

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// PatchUser 	godoc
// @Summary 	Изменение пользователя
// @Description Изменение информации о пользователе в системе по его ID
// @Description ### Логика обновления полей (Three-state logic):
// @Description - Если поле **не передано** в запросе, то оно не изменяется.
// @Description - Если поле **передано со значением null**, то оно обнуляется.
// @Description - Если поле **передано с конкретным значением**, то оно обновляется на это значение.
// @Description Ограничения: `full_name` не может быть null, `phone_number` должен быть в формате e164.
// @Tags 		users
// @Accept 		json
// @Produce 	json
// @Param 		id path int true "ID пользователя"
// @Param 		request body PatchUserRequest true "Новые данные пользователя"
// @Success 	200 {object} PatchUserResponse "Успешное изменение данных пользователя"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 	409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

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
