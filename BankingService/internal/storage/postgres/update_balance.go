package postgres

import (
	"context"
	"errors"
	"fmt"
)

func (p *PostgresRepository) UpdateAccountBalance(ctx context.Context, accountID int64, amount int64) error {
	query := `
		UPDATE accounts 
			SET balance = balance + $1
			WHERE id = $2
	`
	result, err := p.pool.Exec(ctx, query, amount, accountID)
	if err != nil {
		return fmt.Errorf("UpdateAccountBalance: %w", err)
	}
	rows := result.RowsAffected()
	if rows == 0 {
		return errors.New("account not found")
	}
	var transactionType string
	switch {
	case amount > 0:
		transactionType = "deposit"
	default:
		transactionType = "withdraw"
	}
	query = `
		INSERT 
			INTO transactions (account_id, amount, currency, type, status)
			VALUES($1, $2, $3, $4, $5)
	`
	_, err = p.pool.Exec(ctx, query, &accountID, &amount, "RUB", &transactionType, "success")
	return err
}
