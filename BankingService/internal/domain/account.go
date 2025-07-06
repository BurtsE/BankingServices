package domain

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

// Account — банковский счет пользователя
type Account struct {
	ID        int64           `json:"id"`
	UserID    string          `json:"user_id"`
	Number    string          `json:"number"`
	Currency  string          `json:"currency"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"created_at"`
	IsActive  bool            `json:"is_active"`
}

func (a Account) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (a Account) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &a)
}
