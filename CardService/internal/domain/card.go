package domain

import "time"

type Card struct {
	ExpiryMonth    int       `json:"expiry_month"`
	ExpiryYear     int       `json:"expiry_year"`
	ID             int64     `json:"id"`
	AccountID      string    `json:"account_id"`
	PAN            string    `json:"number"` // encode before saving
	CVV            string    `json:"cvv"`    // do not save to storage
	CardholderName string    `json:"cardholder_name"`
	CreatedAt      time.Time `json:"created_at"`
	IsActive       bool      `json:"is_active"`
}

func (c *Card) GenerateTimeExpiry() {
	c.ExpiryYear = time.Now().AddDate(4, 0, 0).Year()
	c.ExpiryMonth = int(time.Now().Month())
}
