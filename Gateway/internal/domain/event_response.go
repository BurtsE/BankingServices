package domain

import "github.com/google/uuid"

type EventResponse struct {
	id   uuid.UUID
	data []byte
	EventStatus
}

func (e *EventResponse) ID() uuid.UUID {
	return e.id
}

func (e *EventResponse) Data() []byte {
	return e.data
}
