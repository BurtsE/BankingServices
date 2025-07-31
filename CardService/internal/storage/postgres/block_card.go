package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// BlockCard finds cards with same pan hash, decodes pan, then changes status
func (p *PostgresRepository) BlockCard(ctx context.Context, accountID string, pan string) (bool, error) {
	var (
		id           int64
		encryptedPan []byte
	)
	query := `
		SELECT id, encode(encrypted_pan, 'escape')::text
			FROM cards
			WHERE account_id = $1 AND hashed_pan = digest( $2, 'sha256')
	`

	rows, err := p.pool.Query(ctx, query, &accountID, &pan)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &encryptedPan)
		if err != nil {
			return false, err
		}

		if pan == string(encryptedPan) {
			break
		}
	}
	var is_active bool
	query = `
		UPDATE cards
			SET is_active = false
			WHERE id = $1
		RETURNING is_active
	`
	err = p.pool.QueryRow(ctx, query, id).Scan(&is_active)
	if err != nil {
		return false, fmt.Errorf("BlockCard: %w", err)
	}
	if is_active {
		return false, nil
	}

	return true, nil
}
