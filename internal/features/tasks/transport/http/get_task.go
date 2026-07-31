package tasks_transport_http

import (
	"net/http"

	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetTaskResponse TaskDTOResponse

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
