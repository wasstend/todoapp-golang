package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(ctx context.Context, queryParams domain.StatisticsFilters) (domain.Statistics, error) {
	if queryParams.From != nil && queryParams.To != nil {
		if queryParams.To.Before(*queryParams.From) || queryParams.To.Equal(*queryParams.From) {
			return domain.Statistics{}, fmt.Errorf(
				"'to' must be after 'from': %w",
				core_errors.ErrorInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, queryParams)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"get tasks from repository: %w",
			err,
		)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil

}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{}
	}

	tasksCreated := len(tasks)
	tasksCompleted := 0

	var totalCompletionDuration time.Duration

	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.CompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	tasksCompletionRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAvgTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)

		tasksAvgTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&tasksCompletionRate,
		tasksAvgTime,
	)
}
