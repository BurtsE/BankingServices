package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventRequest struct {
	id        uuid.UUID
	createdAt time.Time
	data      []byte
	EventType
}

func (e *EventRequest) ID() uuid.UUID {
	return e.id
}

func (e *EventRequest) CreatedAt() time.Time {
	return e.createdAt
}

func (e *EventRequest) Data() []byte {
	return e.data
}

func NewEventRequest(data []byte, eventType string) (EventRequest, error) {
	e := EventRequest{}
	e.id = uuid.New()
	e.createdAt = time.Now()
	e.data = data
	switch eventType {
	case "create_user":
		e.EventType = CreateUserEvent
	default:
		return EventRequest{}, ErrInvalidEventType
	}
	return e, nil
}
