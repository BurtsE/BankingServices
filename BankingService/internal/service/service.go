package service

import (
	"BankingService/internal/domain"
	"BankingService/internal/storage"
	"context"
)

type BankingService interface {
	CreateAccount(ctx context.Context, userID string, currency string) (*domain.Account, error)
	Deposit(ctx context.Context, accountID int64, amount string) error
	Withdraw(ctx context.Context, accountID int64, amount string) error
	Transfer(ctx context.Context, fromAccountID, toAccountID int64, amount string) error
	//GetAccountsByUser(ctx context.Context, userID string) ([]*domain.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error)
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

func (s *Service) GetAccountByID(ctx context.Context, accountID int64) (*domain.Account, error) {
	return s.storage.GetAccountByID(ctx, accountID)
}
