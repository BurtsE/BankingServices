package service

import (
	"BankingService/internal/domain"
	"BankingService/internal/storage"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type BankingService interface {
	CreateAccount(ctx context.Context, userID string,
		currencyStr, accountTypeStr, accountSubTypeStr string) (*domain.Account, error)

	Deposit(ctx context.Context, accountID string, amount string) error

	Withdraw(ctx context.Context, accountID string, amount string) error

	Transfer(ctx context.Context, fromAccountID, toAccountID string, amount string) error

	//GetAccountsByUser(ctx context.Context, userID string) ([]*domain.Account, error)

	GetAccountByID(ctx context.Context, accountID string) (*domain.Account, error)
}

var _ BankingService = (*Service)(nil)

type Service struct {
	storage storage.BankingStorage
}

func NewBankingService(storage storage.BankingStorage) *Service {
	return &Service{storage: storage}
}

func (s *Service) GetAccountsByUser(ctx context.Context, userID string) ([]*domain.Account, error) {
	return s.storage.GetAccountsByUser(ctx, userID)
}

func (s *Service) GetAccountByID(ctx context.Context, accountID string) (*domain.Account, error) {
	accountUUID, err := uuid.Parse(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse account ID %s: %w", accountID, err)
	}

	return s.storage.GetAccountByID(ctx, accountUUID)
}
