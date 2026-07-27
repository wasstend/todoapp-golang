package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/wasstend/todoapp-golang/internal/core/domain"
)

func (r *UserRepository) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
	SELECT * FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query get users error: %w", err)
	}

	defer rows.Close()

	var userModels []UserModel

	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("scan users error: %w", err)
		}

		userModels = append(userModels, userModel)

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil
}
