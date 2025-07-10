package postgres

import (
	"BankingService/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

func (p *PostgresRepository) GetAccountByID(ctx context.Context, accountID uuid.UUID) (*domain.Account, error) {
	query := `
		SELECT id, user_id, number, currency, balance, created_at, is_active
			FROM accounts
			WHERE uuid=$1
	`
	var currencyStr string
	acc := domain.Account{UUID: accountID}
	err := p.pool.QueryRow(ctx, query, accountID.String()).Scan(
		&acc.ID,
		&acc.UserID,
		&acc.Number,
		&currencyStr,
		&acc.Balance,
		&acc.CreatedAt,
		&acc.IsActive,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("GetAccountByID: %w", err)
	}

	currency, err := currencyStringToDomain(currencyStr)
	if err != nil {
		return nil, fmt.Errorf("GetAccountByID: %w", err)
	}

	acc.Currency = currency
	return &acc, nil
}

func (p *PostgresRepository) GetAccountsByUser(ctx context.Context, userID string) ([]*domain.Account, error) {
	query := `
		SELECT id, uuid, number, currency, balance, created_at, is_active
			FROM accounts
			WHERE user_id=$1
	`
	rows, err := p.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("GetAccountsByUser: %w", err)
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var (
			acc         domain.Account
			currencyStr string
			uuidStr     string
		)
		err = rows.Scan(
			&acc.ID,
			&uuidStr,
			&acc.Number,
			&currencyStr,
			&acc.Balance,
			&acc.CreatedAt,
			&acc.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAccountsByUser scan: %w", err)
		}

		UUID, err := uuid.Parse(uuidStr)
		if err != nil {
			return nil, fmt.Errorf("GetAccountsByUser uuid parce: %w", err)
		}
		acc.UUID = UUID

		currency, err := currencyStringToDomain(currencyStr)
		if err != nil {
			return nil, fmt.Errorf("GetAccountByID: %w", err)
		}
		acc.Currency = currency

		accounts = append(accounts, &acc)
	}

	return accounts, rows.Err()
}
