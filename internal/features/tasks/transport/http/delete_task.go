package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

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
