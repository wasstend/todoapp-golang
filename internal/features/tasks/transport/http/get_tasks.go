package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetTasksResponse []TaskDTOResponse

func (h *TasksHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseWriter := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("Invoke GetTasks handler")

	limit, offset, user_ID, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseWriter.ErrorResponse(
			err,
			"failed to get query params",
		)
		return
	}

	tasks, err := h.tasksService.GetTasks(ctx, limit, offset, user_ID)
	if err != nil {
		responseWriter.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(tasksDtoFromDomain(tasks))

	responseWriter.JsonResponse(response, http.StatusOK)

}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIdQueryParamKey = "author_user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIdQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query parameter: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query parameter: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query parameter: %w", err)
	}

	return limit, offset, userID, nil
}
