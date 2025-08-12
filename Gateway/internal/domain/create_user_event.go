package domain

var CreateUserEvent = &createUserRequestEvent{}

type createUserRequestEvent struct {
}

func (c *createUserRequestEvent) isEventType() {}

func (c *createUserRequestEvent) Type() string {
	return "create_user"
}
