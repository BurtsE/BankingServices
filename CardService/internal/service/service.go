package service

import (
	"CardService/internal/domain"
	"context"
)

type CardService interface {
	GenerateVirtualCard(ctx context.Context, accountID string, cardholderName string) (*domain.Card, error)
	GetCardsByAccount(ctx context.Context, accountID string) ([]*domain.Card, error)
	//GetCardByIDForOwner(ctx context.Context, cardID, ownerUserID int64) (*domain.Card, error) // с расшифровкой
}

type IBankingService interface {
	AccountIsActive(ctx context.Context, accountID string) (ok bool, err error)
}
