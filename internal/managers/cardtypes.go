package managers

type CodeType string

const (
	Ean13 CodeType = "ean13"
)

type CardType struct {
	PublicId string
	Name     string
	Code     CodeType
}

var CardTypes = []CardType{
	{
		PublicId: "s7lJTYHX",
		Name:     "Test card",
		Code:     Ean13,
	},
}
