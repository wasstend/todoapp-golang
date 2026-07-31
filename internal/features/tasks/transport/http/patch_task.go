package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_types "github.com/wasstend/todoapp-golang/internal/core/transport/http/types"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate(requestValidator *validator.Validate) error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("Title can't be NULL")
		}

		TitleLen := len([]rune(*r.Title.Value))
		if TitleLen < 1 || TitleLen > 100 {
			return fmt.Errorf("Title length must be minimum 1 and maximum 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			DescriptionLen := len([]rune(*r.Description.Value))
			if DescriptionLen < 1 || DescriptionLen > 1000 {
				return fmt.Errorf("Description length must be minimum 1 and maximum 1000 symbols")
			}
		}

	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("Completed can't be NULL")
		}
	}

	return nil
}

func (h *TasksHTTPHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke PatchTask handler")

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task id",
		)
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskPatch := taskPatchFromRequest(request)
	patchedTask, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchTaskResponse(taskDTOFromDomain(patchedTask))
	responseHandler.JsonResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
