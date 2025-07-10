package domain

var (
	Financial     = &financial{}
	Commercial    = &commercial{}
	NonCommercial = &nonCommercial{}
	Election      = &election{}
	Cooperative   = &cooperative{}
	Defence       = &defense{}
	SpecialBroker = &specialBroker{}
	Physical      = &physical{}
)

type AccountSubType interface {
	isClassification()
	String() string
	Code() string
}

type financial struct{ AccountSubType }

func (f *financial) String() string {
	return "natural"
}

func (f *financial) Code() string {
	return "01"
}

type commercial struct{ AccountSubType }

func (c *commercial) String() string {
	return "commercial"
}

func (c *commercial) Code() string {
	return "01"
}

// non commercial organization
type nonCommercial struct{ AccountSubType }

func (n *nonCommercial) String() string {
	return "nonCommercial"
}

func (n *nonCommercial) Code() string {
	return "03"
}

// election
type election struct{ AccountSubType }

func (e *election) String() string {
	return "election"
}

func (e *election) Code() string {
	return "04"
}

// TFA
type cooperative struct{ AccountSubType }

func (co *cooperative) String() string {
	return "tfa"
}

func (co *cooperative) Code() string {
	return "05"
}

// The head contractor of the defense order
type defense struct{ AccountSubType }

func (d *defense) String() string {
	return "defence"
}

func (d *defense) Code() string {
	return "06"
}

// Special brokerage accounts of type "C"
type specialBroker struct{ AccountSubType }

func (s *specialBroker) String() string {
	return "specialBroker"
}

func (s *specialBroker) Code() string {
	return "07"
}

// physical accounts (only residents)

type physical struct{ AccountSubType }

func (p *physical) String() string {
	return "specialBroker"
}

func (p *physical) Code() string {
	return "17"
}
