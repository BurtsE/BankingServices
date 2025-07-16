package service

import (
	"context"
	"github.com/google/uuid"
)

func (s *Service) AccountIsActive(ctx context.Context, accountID string) (bool, error) {
	id, err := uuid.Parse(accountID)
	if err != nil {
		return false, err
	}

	account, err := s.storage.GetAccountByID(ctx, id)
	if err != nil {
		return false, err
	}

	return account.IsActive, nil
}
