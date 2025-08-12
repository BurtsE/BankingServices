package domain

import "github.com/google/uuid"

type EventResponse struct {
	id uuid.UUID
	EventStatus
}

func (e *EventResponse) ID() uuid.UUID {
	return e.id
}
