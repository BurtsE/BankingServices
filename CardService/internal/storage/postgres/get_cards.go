package postgres

import (
	"CardService/internal/domain"
	"context"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
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
			card          domain.Card
			encrypted_pan []byte
		)

		if err := rows.Scan(&card.ID, &encrypted_pan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt); err != nil {
			return nil, fmt.Errorf("GetCardsByAccount scan: %w", err)
		}

		pan, err := p.decryptWithPGP(encrypted_pan)
		if err != nil {
			return nil, fmt.Errorf("GetCardsByAccount decryption: %w", err)
		}

		card.PAN = pan
		card.AccountID = accountID
		cards = append(cards, &card)
	}

	return cards, rows.Err()
}

// Decrypt armored encrypted message using the private key and obtain the plaintext
func (p *PostgresRepository) decryptWithPGP(encrypted []byte) (string, error) {
	pgp := crypto.PGP()

	// Decrypt armored encrypted message using the private key and obtain the plaintext
	decHandle, err := pgp.Decryption().DecryptionKey(p.privateKey).New()
	if err != nil {
		return "", err
	}

	decrypted, err := decHandle.Decrypt(encrypted, crypto.Armor)
	if err != nil {
		return "", err
	}

	return decrypted.String(), nil
}
