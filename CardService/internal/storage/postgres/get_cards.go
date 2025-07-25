package postgres

import (
	"CardService/internal/domain"
	"context"
	"fmt"
)

func (p *PostgresRepository) GetCardsByAccount(ctx context.Context, accountID string) ([]*domain.Card, error) {
	query := `
		SELECT id, encode(encrypted_pan, 'escape')::text, expiry_month, expiry_year , cardholder_name, is_active, created_at
			FROM cards 
			WHERE account_id=$1
	`
	rows, err := p.pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("GetCardsByAccount: %w", err)
	}
	defer rows.Close()

	cards := make([]*domain.Card, 0)
	for rows.Next() {
		var (
			card         domain.Card
			encryptedPan []byte
		)

		if err = rows.Scan(&card.ID, &encryptedPan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCardsByAccount scan: %w", err)
		}

		pan, err := p.decryptWithPGP(encryptedPan)
		if err != nil {
			return nil, fmt.Errorf("GetCardsByAccount decryption: %w", err)
		}

		card.PAN = pan
		card.AccountID = accountID
		cards = append(cards, &card)
	}

	return cards, rows.Err()
}
