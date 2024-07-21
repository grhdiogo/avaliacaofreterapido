package quote

type Address struct {
	Cep string
}

type Volume struct {
	Category      int
	Amount        int
	UnitaryWeight float64
	Price         float64
	Sku           string
	Height        float64
	Width         float64
	Length        float64
}

type Entity struct {
	ID          string
	CpfCnpj     string
	Address     Address
	Volumes     []Volume
	RawRequest  []byte
	RawResponse []byte
}
