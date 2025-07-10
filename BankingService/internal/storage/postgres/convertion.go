package postgres

import (
	"BankingService/internal/domain"
	"fmt"
)

func currencyStringToDomain(currency string) (domain.Currency, error) {
	switch currency {
	case "USD":
		return domain.USD, nil
	case "RUB":
		return domain.RUB, nil
	default:
		return nil, fmt.Errorf("unknown currency: %s", currency)
	}
}
