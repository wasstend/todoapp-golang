package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/wasstend/todoapp-golang/internal/core/logger"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

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
	limit, err := core_http_utils.GetIntQueryParams(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query parameter: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParams(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query parameter: %w", err)
	}

	return limit, offset, nil
}
