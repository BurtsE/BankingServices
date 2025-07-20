package postgres

import (
	"CardService/internal/domain"
	"context"
	"fmt"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

func (p *PostgresRepository) CreateVirtualCard(ctx context.Context, card *domain.Card) (int64, error) {
	encryptedPan, err := p.encryptWithPGP([]byte(card.PAN))
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	query := `
		INSERT
			INTO cards (account_id, encrypted_pan, expiry_month, expiry_year , cardholder_name, is_active, created_at)
			VALUES ($1, $2::bytea, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id int64
	err = p.pool.QueryRow(ctx, query, &card.AccountID, &encryptedPan, &card.ExpiryMonth, &card.ExpiryYear, &card.CardholderName, &card.IsActive, &card.CreatedAt).Scan(&id)
	return id, err

}

func (p *PostgresRepository) encryptWithPGP(data []byte) (string, error) {
	pgp := crypto.PGP()

	// Encrypt plaintext message using a public key
	encHandle, err := pgp.Encryption().Recipient(p.publicKey).New()
	if err != nil {
		return "", err
	}

	pgpMessage, err := encHandle.Encrypt(data)
	if err != nil {
		return "", err
	}

	armored, err := pgpMessage.Armor()
	if err != nil {
		return "", err
	}

	return armored, nil
}
