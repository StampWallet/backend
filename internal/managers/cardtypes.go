package managers

type CodeType string

const (
	Ean13 CodeType = "ean13"
	Qr    CodeType = "qr"
)

type CardType struct {
	PublicId string
	Name     string
	Code     CodeType
	ImageUrl string
}

var CardTypes = []CardType{
	{
		PublicId: "biedronka",
		Name:     "Moja Biedronka",
		Code:     Ean13,
		ImageUrl: "biedronka.png",
	},
	{
		PublicId: "kaufland",
		Name:     "Kaufland Card",
		Code:     Qr,
		ImageUrl: "kaufland.png",
	},
}

func InitUrls(baseUrl string) {
	for i := 0; i != len(CardTypes); i++ {
		CardTypes[i].ImageUrl = baseUrl + CardTypes[i].ImageUrl
	}
}
