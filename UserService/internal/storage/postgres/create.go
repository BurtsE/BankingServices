package postgres

import (
	model "UserService/internal/domain"
	"context"
)

func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (uuid, email, password, full_name, created_at)
		VALUES ($1, $2, $3, $4, $5)

	`
	_, err := r.pool.Exec(ctx, query, user.UUID, user.Email, user.PasswordHash, user.FullName, user.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
