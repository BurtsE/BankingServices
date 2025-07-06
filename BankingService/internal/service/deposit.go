package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

func (s *Service) Deposit(ctx context.Context, accountID int64, amount string) error {

	value, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if value.IsNegative() {
		return errors.New("amount must be positive")
	}

	return s.storage.UpdateAccountBalance(ctx, accountID, value)
}
