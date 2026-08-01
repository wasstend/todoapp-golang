package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetTasksResponse []TaskDTOResponse

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// GetTasks 	godoc
// @Summary 	Получение задач
// @Description Получение списка задач с фильтрацией по автору и/или опциональной пагинацией
// @Tags 		tasks
// @Produce 	json
// @Param 		author_user_id query int false "ID автора задач"
// @Param 		limit query int false "Размер страницы с задачами"
// @Param 		offset query int false "Смещение страницы с задачами"
// @Success 	200 {object} GetTasksResponse "Успешное получение задач"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks [get]
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
