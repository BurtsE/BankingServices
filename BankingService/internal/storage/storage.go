package storage

import (
	model "BankingService/internal/domain"
	"context"
	"github.com/shopspring/decimal"
)

type BankingStorage interface {
	BeginTransaction(ctx context.Context) (Transaction, error)
	CreateAccount(ctx context.Context, userID string, currency string) (*model.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*model.Account, error)
	GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error)
	UpdateAccountBalance(ctx context.Context, accountID int64, amount decimal.Decimal) error
}

type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
