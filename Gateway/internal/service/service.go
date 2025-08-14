package service

import (
	"context"
	"gateway/internal/domain"
)

type IUserService interface {
	Validate(jwtToken string) (uuid string, err error)
}

// Should be thread safe
type EventService interface {
	AddEvent(context.Context, domain.EventRequest) error
}
