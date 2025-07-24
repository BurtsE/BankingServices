package card_service

import (
	"CardService/internal/domain"
	"context"
	"fmt"
)

func (c *CardService) GetCardsByAccount(ctx context.Context, accountID string) ([]*domain.Card, error) {
	var err error
	cards, err := c.storage.GetCardsByAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("get cards by account: %w", err)
	}
	for i := range cards {
		cards[i].CVV, err = c.generateCVV(cards[i].PAN, cards[i].ExpiryYear, cards[i].ExpiryMonth)
		if err != nil {
			return nil, fmt.Errorf("generate cvv: %w", err)
		}
	}
	return cards, nil
}
