package card_service

import (
	"CardService/internal/config"
	"CardService/internal/domain"
	"CardService/internal/service"
	"CardService/internal/storage"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"log"
)

var _ service.CardService = (*CardService)(nil)

type CardService struct {
	storage storage.CardStorage
	gcm     cipher.AEAD
}

func NewCardService(storage storage.CardStorage) *CardService {
	s := &CardService{storage: storage}

	// generate a new aes cipher using 32 byte long key
	c, err := aes.NewCipher([]byte(config.GetCardSecretKey()))
	if err != nil {
		log.Fatalf("%v", err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("%v", err)
	}
	s.gcm = gcm

	return s
}

func (c *CardService) GetCardsByAccount(ctx context.Context, accountID int64) ([]*domain.Card, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CardService) GetCardByIDForOwner(ctx context.Context, cardID, ownerUserID int64) (*domain.Card, error) {
	//TODO implement me
	panic("implement me")
}
