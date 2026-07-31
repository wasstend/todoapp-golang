package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
	core_postgres_pool "github.com/wasstend/todoapp-golang/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(ctx context.Context, id int) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
	SELECT * FROM todoapp.tasks
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var taskModel TaskModel
	if err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d': %w",
				id,
				core_errors.ErrorNotFound,
			)
		}
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
