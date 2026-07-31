package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
	core_postgres_pool "github.com/wasstend/todoapp-golang/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(ctx context.Context, id int, task domain.Task) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.tasks
	SET
		title=$1,
		description=$2,
		completed=$3,
		completed_at=$4,
		version=version+1

	WHERE id=$5 AND version=$6

	RETURNING *
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		task.ID,
		task.Version,
	)

	var taskModel TaskModel
	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d' concurrently accessed: %w",
				id,
				core_errors.ErrorConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("scan error: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil

}
