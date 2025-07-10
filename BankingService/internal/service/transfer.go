package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Service) Transfer(ctx context.Context, fromAccountID, toAccountID string, amount string) error {

	if fromAccountID == toAccountID {
		return errors.New("cannot transfer to the same account")
	}

	value, err := decimal.NewFromString(amount)
	if err != nil {
		return fmt.Errorf("invalid amount: %w", err)
	}

	if value.IsNegative() {
		return errors.New("amount must be positive")
	}

	fromAccountUUID, err := uuid.Parse(fromAccountID)
	if err != nil {
		return fmt.Errorf("could not parse account ID %s: %w", fromAccountID, err)
	}

	from, err := s.storage.GetAccountByID(ctx, fromAccountUUID)
	if err != nil {
		return err
	}

	toAccountUUID, err := uuid.Parse(toAccountID)
	if err != nil {
		return fmt.Errorf("could not parse account ID %s: %w", toAccountID, err)
	}

	to, err := s.storage.GetAccountByID(ctx, toAccountUUID)
	if err != nil {
		return err
	}

	if from.Currency != to.Currency {
		return errors.New("insufficient currencies")
	}

	if from.Balance.LessThan(value) {
		return errors.New("insufficient balance")
	}

	tx, err := s.storage.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err = s.storage.UpdateAccountBalance(ctx, fromAccountID, value.Mul(decimal.NewFromInt(-1))); err != nil {
		return err
	}

	if err = s.storage.UpdateAccountBalance(ctx, toAccountID, value); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}
