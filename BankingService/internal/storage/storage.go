package storage

import (
	model "BankingService/internal/domain"
	"context"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type BankingStorage interface {
	BeginTransaction(ctx context.Context) (Transaction, error)

	CreateAccount(ctx context.Context, account *model.Account) (int64, error)

	GetAccountsNumber(ctx context.Context) (int64, error)

	GetAccountByID(ctx context.Context, accountID uuid.UUID) (*model.Account, error)

	GetAccountsByUser(ctx context.Context, userID string) ([]*model.Account, error)

	UpdateAccountBalance(ctx context.Context, accountID string, amount decimal.Decimal) error
}

type Transaction interface {
	Commit(context.Context) error

	Rollback(context.Context) error
}
