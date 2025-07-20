package card_service

import (
	"CardService/internal/domain"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"
)

func (c *CardService) GenerateVirtualCard(ctx context.Context, accountID string, cardholderName string) (*domain.Card, error) {
	card := &domain.Card{}
	card.GenerateTimeExpiry()

	pan, err := c.generatePAN()
	if err != nil {
		return nil, err
	}

	cvv, err := c.generateCVV(pan, card.ExpiryYear, card.ExpiryMonth)
	if err != nil {
		return nil, err
	}

	card.AccountID = accountID
	card.PAN = pan
	card.CVV = cvv
	card.CardholderName = cardholderName
	card.CreatedAt = time.Now()
	card.IsActive = true

	id, err := c.storage.CreateVirtualCard(ctx, card)
	if err != nil {
		return nil, err
	}
	card.ID = id

	return card, nil
}

func (c *CardService) generateCVV(pan string, expiryYear, expiryMonth int) (string, error) {
	// Create a hash of the combined input
	combined := pan + "|" + time.Date(expiryYear, time.Month(expiryMonth), 0, 0, 0, 0, 0, time.Local).String()

	// Hash the combined input to get a consistent key
	hasher := sha256.New()
	hasher.Write([]byte(combined))
	key := hasher.Sum(nil)[:aes.BlockSize] // Use first 16 bytes as AES key

	// Create a cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// We'll use a fixed nonce (since input is always the same)
	nonce := make([]byte, gcm.NonceSize())

	// Encrypt some fixed data (we just need consistent output)
	// Using the variables again as plaintext to incorporate all input
	ciphertext := gcm.Seal(nil, nonce, []byte(combined), nil)

	// Convert first 4 bytes of ciphertext to a number
	num := binary.BigEndian.Uint32(ciphertext[:4])

	// Reduce to three digits (000-999)
	cvv := make([]byte, 0, 3)
	for range 3 {
		cvv = append(cvv, byte(num%10)+48)
		num /= 10
	}
	return string(cvv), nil
}

func (c *CardService) generatePAN() (string, error) {
	pan := make([]byte, 15)
	pan[0] = byte(4)
	sum := 8

	for i := range 14 {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generating card PAN: %v", err)
		}

		pan[i+1] = byte(digit.Int64())
		if i%2 != 0 {
			sum += int(digit.Int64()) * 2
		} else {
			sum += int(digit.Int64())
		}
	}

	lastDigit := 0
	for (sum+lastDigit)%10 != 0 {
		lastDigit++
	}
	pan = append(pan, byte(lastDigit))
	for i := range pan {
		pan[i] += 48
	}

	return string(pan), nil
}
