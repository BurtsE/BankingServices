package card_service

import "context"

func (c *CardService) BlockCard(ctx context.Context, accountID string, pan string) (bool, error) {
	blocked, err := c.storage.BlockCard(ctx, accountID, pan)
	if err != nil {
		return false, err
	}

	return blocked, nil
}
