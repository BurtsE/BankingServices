package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventRequest struct {
	id        uuid.UUID
	createdAt time.Time
	EventStatus
	EventType
}

func (e *EventRequest) ID() uuid.UUID {
	return e.id
}

func NewEventRequest(value string) (EventRequest, error) {
	e := EventRequest{}
	e.id = uuid.New()
	switch value {
	case "create_user":
		e.EventType = CreateUserEvent
	default:
		return EventRequest{}, ErrInvalidEventType
	}
	return e, nil
}
