package postgres

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
)

func (p *PostgresRepository) UpdateAccountBalance(ctx context.Context, accountID int64, amount decimal.Decimal) error {
	query := `
		UPDATE accounts 
			SET balance = balance + $1
			WHERE id = $2
		RETURNING currency;
	`
	
	var currency string
	err := p.pool.QueryRow(ctx, query, amount, accountID).Scan(&currency)
	if err != nil {
		return fmt.Errorf("UpdateAccountBalance: %w", err)
	}

	var transactionType string
	switch {
	case amount.IsPositive():
		transactionType = "deposit"
	default:
		transactionType = "withdraw"
	}

	query = `
		INSERT 
			INTO transactions (account_id, amount, currency, type, status)
			VALUES($1, $2, $3, $4, $5)
	`
	_, err = p.pool.Exec(ctx, query, &accountID, &amount, &currency, &transactionType, "success")
	return err
}
