package service

import (
	"BankingService/internal/domain"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

func (s *Service) CreateAccount(ctx context.Context, userID string,
	currencyStr, accountTypeStr, accountSubTypeStr string) (*domain.Account, error) {

	currency, err := validateCurrency(currencyStr)
	if err != nil {
		return nil, err
	}

	accountType, err := validateAccountType(accountTypeStr)
	if err != nil {
		return nil, err
	}

	accountSubType, err := validateAccountSubType(accountType, accountSubTypeStr)
	if err != nil {
		return nil, err
	}

	number, err := s.storage.GetAccountsNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts number: %w", err)
	}

	accountNumber := fmt.Sprintf("%s%s%s0000%d",
		accountType.Code(),
		accountSubType.Code(),
		currency.Code(),
		number,
	)

	account := &domain.Account{
		UUID:           uuid.New(),
		UserID:         userID,
		AccountType:    accountType,
		AccountSubType: accountSubType,
		Number:         accountNumber,
		Currency:       currency,
		Balance:        decimal.New(0, 0),
		CreatedAt:      time.Now(),
		IsActive:       true,
	}

	id, err := s.storage.CreateAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	account.ID = id

	return account, nil
}
