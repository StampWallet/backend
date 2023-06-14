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
		ImageUrl: "https://prowly-uploads.s3.eu-west-1.amazonaws.com/uploads/8222/assets/132267/original-190120148b33da12ed35edc531508409.jpg",
	},
	{
		PublicId: "kaufland",
		Name:     "Kaufland Card",
		Code:     Qr,
		ImageUrl: "https://upload.wikimedia.org/wikipedia/commons/f/fc/Kaufland_supermarket2.jpg",
	},
}
