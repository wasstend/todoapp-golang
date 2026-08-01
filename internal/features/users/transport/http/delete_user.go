package users_transport_http

import (
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// DeleteUser 	godoc
// @Summary 	Удаление пользователя
// @Description Удаление пользователя в системе по его ID
// @Tags 		users
// @Param 		id path int true "ID удаляемого пользователя"
// @Success 	204 "Успешное удаление пользователя"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke DeleteUser Handler")

	userId, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id path value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userId); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
