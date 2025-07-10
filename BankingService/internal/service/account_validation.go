package service

import (
	"BankingService/internal/domain"
	"errors"
)

var (
	errInvalidAccountType    = errors.New("invalid account type")
	errInvalidAccountSubType = errors.New("invalid account subtype")
	errInvalidCurrency       = errors.New("invalid currency")
)

func validateCurrency(currencyStr string) (domain.Currency, error) {
	switch currencyStr {
	case "RUB":
		return domain.RUB, nil
	case "USD":
		return domain.USD, nil
	}

	return nil, errInvalidCurrency
}

func validateAccountType(accountType string) (domain.AccountType, error) {
	switch accountType {
	case "natural":
		return domain.NaturalPerson, nil
	case "legal":
		return domain.LegalEntity, nil
	}

	return nil, errors.New("invalid account type")
}

func validateAccountSubType(accountType domain.AccountType, accountSubtype string) (domain.AccountSubType, error) {
	switch accountType {
	case domain.NaturalPerson:
		return validateNaturalPersonSubtype(accountSubtype)
	case domain.LegalEntity:
		return validateLegalEntitySubtype(accountSubtype)

	}

	return nil, errors.New("invalid account type")
}

func validateNaturalPersonSubtype(accountSubtype string) (domain.AccountSubType, error) {
	switch accountSubtype {
	case "physical":
		return domain.Physical, nil
	}
	return nil, errInvalidAccountSubType
}

func validateLegalEntitySubtype(accountSubtype string) (domain.AccountSubType, error) {
	switch accountSubtype {
	case "financial":
		return domain.Financial, nil
	case "commercial":
		return domain.Commercial, nil
	case "nonCommercial":
		return domain.NonCommercial, nil
	case "election":
		return domain.Election, nil
	case "tfa":
		return domain.Cooperative, nil
	case "defence":
		return domain.Defence, nil
	case "specialBroker":
		return domain.SpecialBroker, nil
	}

	return nil, errInvalidAccountSubType
}
