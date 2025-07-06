package postgres

import (
	model "BankingService/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p *PostgresRepository) GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error) {
	query := `
		SELECT id, user_id, currency, balance
			FROM accounts
			WHERE id=$1
	`
	var acc model.Account
	err := p.pool.QueryRow(ctx, query, accountID).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Currency,
		&acc.Balance,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("GetAccountByID: %w", err)
	}
	return &acc, nil
}

func (p *PostgresRepository) GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error) {
	query := `
		SELECT id, user_id, currency, balance
			FROM accounts
			WHERE user_id=$1
	`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAccountsByUser: %w", err)
	}
	defer rows.Close()

	var accounts []*model.Account
	for rows.Next() {
		var acc model.Account
		if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Currency, &acc.Balance); err != nil {
			return nil, fmt.Errorf("GetAccountsByUser scan: %w", err)
		}
		accounts = append(accounts, &acc)
	}
	return accounts, rows.Err()
}
