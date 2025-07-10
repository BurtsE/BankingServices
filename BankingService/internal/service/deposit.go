package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

func (s *Service) Deposit(ctx context.Context, accountID string, amount string) error {

	value, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if value.IsNegative() {
		return errors.New("amount must be positive")
	}

	tx, err := s.storage.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	
	err = s.storage.UpdateAccountBalance(ctx, accountID, value)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
