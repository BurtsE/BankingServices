package storage

import (
	"CardService/internal/domain"
	"context"
)

type CardStorage interface {
	CreateVirtualCard(ctx context.Context, card *domain.Card) (int64, error)
	GetCardsByAccount(ctx context.Context, accountID string) ([]*domain.Card, error)
}
