package quote

type Entity struct {
	Carrier []Carrier `json:"carrier"`
}

type Carrier struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}
