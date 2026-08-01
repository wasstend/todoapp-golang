package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// DeleteTask 	godoc
// @Summary 	Удаление задачи
// @Description Удаление задачи в системе по её ID
// @Tags 		tasks
// @Param 		id path int true "ID удаляемой задачи"
// @Success 	204 "Успешное удаление задачи"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not Found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [delete]
func (h *TasksHTTPHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke DeleteTask handler")

	id, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task id",
		)
		return
	}

	if err := h.tasksService.DeleteTask(ctx, id); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete task",
		)
		return
	}

	responseHandler.NoContentResponse()
}
