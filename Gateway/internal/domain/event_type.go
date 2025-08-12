package domain

// Enums for event type
type EventType interface {
	isEventType()
	Type() string
}
