package users_transport_http

import (
	"fmt"
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// GetUsers 	godoc
// @Summary 	Получение пользователей
// @Description Получение данных о всех пользователях в системе с опциональной пагинацией
// @Tags 		users
// @Produce 	json
// @Param 		limit query int false "Размер страницы с пользователями"
// @Param 		offset query int false "Смещение страницы с пользователями"
// @Success 	200 {object} GetUsersResponse "Успешное получение пользователей"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/users [get]
func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke GetUsers Handler")

	limit, offset, err := getLimitOffset(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit and offset query parameters",
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
		)
		return
	}

	response := GetUsersResponse(usersDTOFromDomains(userDomains))

	responseHandler.JsonResponse(response, http.StatusOK)

}

func getLimitOffset(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query parameter: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query parameter: %w", err)
	}

	return limit, offset, nil
}
