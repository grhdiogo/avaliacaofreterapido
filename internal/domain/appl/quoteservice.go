package appl

import "avaliacaofreterapido/internal/domain/quote"

type QuoteRequestParams struct {
	Recipient Recipient `json:"recipient"`
	Volumes   []Volume  `json:"volumes"`
}

type Recipient struct {
	Address Address `json:"address"`
}

type Address struct {
	Zipcode string `json:"zipcode"`
}

type Volume struct {
	Category      int     `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight int     `json:"unitary_weight"`
	Price         int     `json:"price"`
	Sku           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

type QuoteService interface {
	GetQuotes(r QuoteRequestParams) (*quote.Entity, error)
}

type quoteService struct{}

func (s *quoteService) GetQuotes(r QuoteRequestParams) (*quote.Entity, error) {
	return &quote.Entity{}, nil
}

func NewQuoteService() QuoteService {
	return &quoteService{}
}
