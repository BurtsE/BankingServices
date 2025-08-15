package domain


var DepositAccountEvent = &depositAccountRequestEvent{}

type depositAccountRequestEvent struct {
}

func (c *depositAccountRequestEvent) isEventType() {}

func (c *depositAccountRequestEvent) Type() string {
	return "deposit_account"
}