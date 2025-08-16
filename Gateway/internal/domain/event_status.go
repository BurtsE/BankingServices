package domain

// Enums for event status
var (
	Succeess   = &success{}
	Fail       = &fail{}
	InProgress = &inProgress{}
)

type EventStatus interface {
	isEventStatus()
	Status() string
}

type inProgress struct {
}

func (i *inProgress) isEventStatus() {}

func (i *inProgress) Status() string {
	return "in_progress"
}

type success struct {
}

func (s *success) isEventStatus() {}

func (s *success) Status() string {
	return "succeess"
}

type fail struct {
}

func (f *fail) isEventStatus() {}

func (f *fail) Status() string {
	return "fail"
}
