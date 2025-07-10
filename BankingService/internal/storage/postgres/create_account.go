package postgres

import (
	"BankingService/internal/domain"
	"context"
	"fmt"
)

func (p *PostgresRepository) CreateAccount(ctx context.Context, account *domain.Account) (int64, error) {
	query := `
		INSERT INTO 
			accounts (uuid, user_id, number, currency, balance, created_at, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	currency := account.Currency.String()
	var id int64
	err := p.pool.QueryRow(ctx, query,
		&account.UUID,
		&account.UserID,
		&account.Number,
		&currency,
		&account.Balance,
		&account.CreatedAt,
		&account.IsActive,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateAccount: %w", err)
	}

	return id, nil
}
