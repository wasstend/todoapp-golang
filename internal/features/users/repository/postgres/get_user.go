package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/wasstend/todoapp-golang/internal/core/domain"
	core_errors "github.com/wasstend/todoapp-golang/internal/core/errors"
)

func (r *UserRepository) GetUser(ctx context.Context, id int) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetTimeout())
	defer cancel()

	query := `
	SELECT * FROM todoapp.users
	WHERE id=$1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d': %w",
				id,
				core_errors.ErrorNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
