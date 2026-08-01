package users_transport_http

import (
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// GetUser 	godoc
// @Summary 	Получить пользователя
// @Description Получение пользователя в системе по его ID
// @Tags 		users
// @Produce 	json
// @Param 		id path int true "ID пользователя"
// @Success 	200 {object} GetUserResponse "Успешное получение данных о пользователе"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [get]
func (h *UsersHTTPHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke GetUser Handler")

	var (
		id = "id"
	)

	userId, err := core_http_request.GetIntPathValue(r, id)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id path value",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userId)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))
	responseHandler.JsonResponse(response, http.StatusOK)
}
