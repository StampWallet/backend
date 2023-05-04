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
		PublicId: "test",
		Name:     "Test card",
		Code:     Ean13,
	},
}
