package postgres

import (
	"CardService/internal/domain"
	"context"
	"crypto/sha256"
	"fmt"
)

func (p *PostgresRepository) CreateVirtualCard(ctx context.Context, card *domain.Card) (int64, error) {
	encryptedPan, err := p.encryptWithPGP([]byte(card.PAN))
	hashedPan := sha256.Sum256([]byte(card.PAN))
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	query := `
		INSERT
			INTO cards (account_id, encrypted_pan, hashed_pan, expiry_month, expiry_year , cardholder_name, is_active, created_at)
			VALUES ($1, $2::bytea, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int64
	err = p.pool.QueryRow(ctx, query, &card.AccountID, &encryptedPan, &hashedPan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt).Scan(&id)
	return id, err

}
