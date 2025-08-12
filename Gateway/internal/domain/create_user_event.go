package domain

var CreateUserEvent = &createUserEvent{}

type createUserEvent struct {
}

func (c *createUserEvent) isEventType() {}

func (c *createUserEvent) Value() string {
	return "create_user"
}
