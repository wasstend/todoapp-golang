package tasks_service

import (
	"context"
	"fmt"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

func (s *TasksService) GetTasks(ctx context.Context, limit, offset, userID *int) ([]domain.Task, error) {

	params := []struct {
		name  string
		value *int
	}{
		{"limit", limit},
		{"offset", offset},
	}

	for _, v := range params {
		if v.value != nil && *v.value < 0 {
			return nil, fmt.Errorf(
				"%s can't be negative: %w",
				v.name,
				core_errors.ErrorInvalidArgument,
			)
		}
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, limit, offset, userID)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get tasks from repository: %w",
			err,
		)
	}

	return tasks, nil
}
