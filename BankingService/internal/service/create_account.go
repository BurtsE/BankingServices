package service

import (
	model "BankingService/internal/domain"
	"context"
)

func (s *Service) CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error) {
	return s.storage.CreateAccount(ctx, userID, currency)
}
