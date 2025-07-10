package domain

// TODO add credit account type, government
var (
	NaturalPerson = &naturalPerson{}
	LegalEntity   = &legalEntity{}
	//Government    = &government{}
)

type AccountType interface {
	isAccountType()
	Code() string
	String() string
}

// NaturalPerson account created for natural person
type naturalPerson struct {
	AccountType
	classification AccountSubType
}

func (p *naturalPerson) String() string {
	return "natural"
}

func (p *naturalPerson) Code() string {
	return "408"
}

// legalEntity account created for legal entity

type legalEntity struct {
	AccountType
}

func (l *legalEntity) String() string {
	return "legal"
}

func (l *legalEntity) Code() string {
	return "407"
}

// government account created for government organization

type government struct {
	AccountType
	classification AccountSubType
}

func (g *government) String() string {
	return "government"
}

func (g *government) Code() string {
	return "406"
}
