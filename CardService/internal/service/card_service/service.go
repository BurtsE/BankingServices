package card_service

import (
	"CardService/internal/config"
	"CardService/internal/service"
	"CardService/internal/storage"
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
