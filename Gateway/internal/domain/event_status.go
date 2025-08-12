package domain

// Enums for event status
var (
	Succeess = &success{}
	Fail     = &fail{}
)

type EventStatus interface {
	isEventStatus()
	Status() string
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
