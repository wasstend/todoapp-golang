package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(ctx context.Context, queryParams domain.StatisticsFilters) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
	SELECT * FROM todoapp.tasks
	`)

	args, conditions := queryParamsToArgs(queryParams)

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(" ORDER BY id ASC")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf(
			"select tasks: %w",
			err,
		)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		err := rows.Scan(
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
			return nil, fmt.Errorf(
				"scan tasks: %w",
				err,
			)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	taskDomains := taskDomainsFromModels(taskModels)

	return taskDomains, nil
}

func queryParamsToArgs(queryParams domain.StatisticsFilters) ([]any, []string) {
	args := []any{}
	conditions := []string{}

	if queryParams.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, queryParams.UserID)
	}

	if queryParams.From != nil {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, queryParams.From)
	}

	if queryParams.To != nil {
		conditions = append(conditions, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, queryParams.To)
	}

	return args, conditions
}
