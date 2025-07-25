package postgres

import "github.com/ProtonMail/gopenpgp/v3/crypto"

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
