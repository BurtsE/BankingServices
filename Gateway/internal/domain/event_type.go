package domain

import "github.com/google/uuid"

type EventType interface {
	isEventType()
	Value() string
}

type Event struct {
	id uuid.UUID
	EventType
}

func (e *Event) ID() uuid.UUID {
	return e.id
}

func NewEvent(value string) (Event, error) {
	e := Event{}
	e.id = uuid.New()
	switch value {
	case "create_user":
		e.EventType = CreateUserEvent
	default:
		return Event{}, ErrInvalidEventType
	}
	return e, nil
}
