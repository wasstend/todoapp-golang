package statistics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_response "github.com/wasstend/todoapp-golang/internal/core/transport/http/response"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"10"`
	TasksCompleted             int      `json:"tasks_completed" example:"5"`
	TasksCompletionRate        *float64 `json:"tasks_completion_rate" example:"33.33"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"2h30m"`
}

// referenced only in swagger annotations below
var _ core_http_response.ErrorResponse

// GetStatistics 	godoc
// @Summary 		Получить статистику
// @Description 	Получение статистики по задачам с опциональной фильтрацией по автору и диапазону дат
// @Tags 			statistics
// @Produce 		json
// @Param 			user_id query int false "ID автора задач"
// @Param 			from 	query string false "Начальная дата (включительно) в формате YYYY-MM-DD"
// @Param 			to 		query string false "Конечная дата (не включительно) в формате YYYY-MM-DD"
// @Success 		200 	{object} GetStatisticsResponse "Успешное получение статистики"
// @Failure 		400 	{object} core_http_response.ErrorResponse "Bad request"
// @Failure 		500 	{object} core_http_response.ErrorResponse "Internal server error"
// @Router 			/statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx, logger, responseHandler := core_http_utils.GetCtxLogResp(w, r)

	logger.Debug("invoke GetStatistics handler")

	queryParams, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get query params",
		)
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, queryParams)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get statistics",
		)
		return
	}

	response := toDTOFromDomain(statistics)

	responseHandler.JsonResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletionRate:        statistics.TasksCompletionRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (domain.StatisticsFilters, error) {
	const (
		layout              = "2006-01-02"
		userIDQueryParamKey = "user_id"
		fromQueryParamKey   = "from"
		toQueryParamKey     = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return domain.StatisticsFilters{}, fmt.Errorf(
			"get 'user_id' query param: %w",
			err,
		)
	}

	from, err := core_http_request.GetTimeQueryParam(r, fromQueryParamKey, layout)
	if err != nil {
		return domain.StatisticsFilters{}, fmt.Errorf(
			"get 'from' query param: %w",
			err,
		)
	}

	to, err := core_http_request.GetTimeQueryParam(r, toQueryParamKey, layout)
	if err != nil {
		return domain.StatisticsFilters{}, fmt.Errorf(
			"get 'to' query param: %w",
			err,
		)
	}

	queryParams := domain.StatisticsFilters{
		UserID: userID,
		From:   from,
		To:     to,
	}

	return queryParams, nil
}
