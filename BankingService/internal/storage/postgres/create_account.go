package postgres

import (
	model "BankingService/internal/domain"
	"context"
	"fmt"
)

func (p *PostgresRepository) CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error) {
	query := `
		INSERT INTO 
			accounts (user_id, currency, balance)
			VALUES ($1, $2, 0.0)
		RETURNING id, user_id, currency, balance, is_active
	`
	var acc model.Account
	err := p.pool.QueryRow(ctx, query, userID, currency).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Currency,
		&acc.Balance,
		&acc.IsActive,
	)
	if err != nil {
		return nil, fmt.Errorf("CreateAccount: %w", err)
	}
	return &acc, nil
}
