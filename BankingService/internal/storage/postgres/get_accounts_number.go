package postgres

import (
	"context"
	"fmt"
)

func (p *PostgresRepository) GetAccountsNumber(ctx context.Context) (int64, error) {
	query := `
		SELECT coalesce(MAX(id), 0)
			FROM accounts
	`
	var id int64
	err := p.pool.QueryRow(ctx, query).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("GetAccountsNumber: %w", err)
	}

	return id, nil
}
