package model

import (
	"encoding/json"
	"time"
)

// User — данные пользователя
type User struct {
	UUID         string    `json:"uuid"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // bcrypt hash
	FullName     string    `json:"full_name"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}
