package domain

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

// Account — user's bank account
type Account struct {
	ID             int64
	UUID           uuid.UUID       `json:"uuid"`
	UserID         string          `json:"user_id"`
	AccountType    AccountType     `json:"-"`
	AccountSubType AccountSubType  `json:"-"`
	Number         string          `json:"number"`
	Currency       Currency        `json:"currency"`
	Balance        decimal.Decimal `json:"balance"`
	CreatedAt      time.Time       `json:"created_at"`
	IsActive       bool            `json:"is_active"`
}

func (a Account) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a Account) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}
