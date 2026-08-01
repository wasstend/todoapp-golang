package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetTaskResponse TaskDTOResponse

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// GetTask 	godoc
// @Summary 	Получить задачу
// @Description Получение задачи в системе по её ID
// @Tags 		tasks
// @Produce 	json
// @Param 		id path int true "ID задачи"
// @Success 	200 {object} GetTaskResponse "Успешное получение данных о задаче"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke GetTask handler")

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id",
		)
		return
	}

	task, err := h.tasksService.GetTask(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)
		return
	}

	response := GetTaskResponse(taskDTOFromDomain(task))

	responseHandler.JsonResponse(response, http.StatusOK)
}
