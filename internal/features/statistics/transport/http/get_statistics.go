package statistics_transport_http

import (
	"fmt"
	"net/http"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_http_request "github.com/wasstend/todoapp-golang/internal/core/transport/http/request"
	core_http_utils "github.com/wasstend/todoapp-golang/internal/core/transport/http/utils"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletionRate        *float64 `json:"tasks_completion_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

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
