package domain

var CreateAccountEvent = &createAccountRequestEvent{}

type createAccountRequestEvent struct {
}

func (c *createAccountRequestEvent) isEventType() {}

func (c *createAccountRequestEvent) Type() string {
	return "create_account"
}