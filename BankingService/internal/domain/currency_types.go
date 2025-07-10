package domain

import (
	"encoding/json"
)

var (
	RUB = &rub{}
	USD = &usd{}
)

type Currency interface {
	isCurrencyType()
	Code() string
	String() string
}

type rub struct {
	Currency
}

func (r *rub) Code() string   { return "810" }
func (r *rub) String() string { return "RUB" }

func (r *rub) MarshalBinary() ([]byte, error) {
	return []byte(r.String()), nil
}
func (r *rub) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

type usd struct {
	Currency
}

func (u *usd) isCurrencyType() {}
func (u *usd) String() string  { return "USD" }
func (u *usd) Code() string    { return "840" }
func (u *usd) MarshalBinary() ([]byte, error) {
	return []byte(u.String()), nil
}
func (u *usd) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}
