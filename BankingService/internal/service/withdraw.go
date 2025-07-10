package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) Withdraw(ctx context.Context, accountID string, amount string) error {

	value, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if value.IsNegative() {
		return errors.New("amount must be positive")
	}

	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		return fmt.Errorf("could not parse account ID %s: %w", accountID, err)
	}

	account, err := s.storage.GetAccountByID(ctx, accountUUID)
	if err != nil {
		return err
	}

	if account.Balance.LessThan(value) {
		return errors.New("insufficient balance")
	}

	return s.storage.UpdateAccountBalance(ctx, accountID, value.Mul(decimal.NewFromInt(-1)))
}
