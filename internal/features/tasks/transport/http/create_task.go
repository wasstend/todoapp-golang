package tasks_transport_http

import (
	"net/http"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"My Task"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Description of my task"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"1"`
}

type CreateTaskResponse TaskDTOResponse

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// CreateTask 	godoc
// @Summary 	Создать задачу
// @Description Создание новой задачи в системе
// @Tags 		tasks
// @Accept 		json
// @Produce 	json
// @Param 		request body CreateTaskRequest true "CreateTask тело запроса"
// @Success 	201 {object} CreateTaskResponse "Успешно созданная задача"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404 {object} core_http_response.ErrorResponse "Not Found (автор задачи не найден)"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks [post]
func (h *TasksHTTPHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke CreateTask Handler")

	var request CreateTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	taskDomain := domainFromDto(request)
	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to create task",
		)
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JsonResponse(response, http.StatusCreated)
}

func domainFromDto(dto CreateTaskRequest) domain.Task {
	return domain.NewTaskUninitialized(
		dto.Title,
		dto.Description,
		dto.AuthorUserID,
	)
}
