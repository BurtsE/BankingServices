package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

func (s *Service) Withdraw(ctx context.Context, accountID int64, amount string) error {

	value, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if value.IsNegative() {
		return errors.New("amount must be positive")
	}

	account, err := s.storage.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	if account.Balance.LessThan(value) {
		return errors.New("insufficient balance")
	}

	return s.storage.UpdateAccountBalance(ctx, accountID, value.Mul(decimal.NewFromInt(-1)))
}
