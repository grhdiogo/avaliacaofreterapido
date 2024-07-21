package resource

import (
	"avaliacaofreterapido/internal/appl"
	"avaliacaofreterapido/internal/interf"
	"encoding/json"
	"net/http"
)

// ========================================================================
// CREATE QUOTATION
// ========================================================================

type CreateQuoteVolumeRequest struct {
	Category      int     `json:"category"`
	Amount        int     `json:"amount"`
	UnitaryWeight float64 `json:"unitary_weight"`
	Price         float64 `json:"price"`
	Sku           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

type CreateQuoteRequest struct {
	Recipient CreateQuoteRequestRecipient `json:"recipient"`
	Volumes   []CreateQuoteVolumeRequest  `json:"volumes"`
}

type CreateQuoteRequestRecipient struct {
	Address CreateQuoteRequestAddress `json:"address"`
}

type CreateQuoteRequestAddress struct {
	Zipcode string `json:"zipcode"`
}

type CreateQuoteCarrierResponse struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}

type CreateQuoteResponse struct {
	Carrier []CreateQuoteCarrierResponse `json:"carrier"`
}

func CreateQuote(r *http.Request) (any, *interf.ErrorHandler) {
	request := new(CreateQuoteRequest)
	// decode body
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(request)
	if err != nil {
		return nil, &interf.ErrorHandler{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			ErrCode:    999999999,
		}
	}
	// service
	app := appl.NewQuoteService(r.Context())
	// volumes
	vlms := make([]appl.Volume, 0)
	for _, v := range request.Volumes {
		vlms = append(vlms, appl.Volume{
			Category:      v.Category,
			Amount:        v.Amount,
			UnitaryWeight: v.UnitaryWeight,
			Price:         v.Price,
			Sku:           v.Sku,
			Height:        v.Height,
			Width:         v.Width,
			Length:        v.Length,
		})
	}
	// make quotation
	resp, err := app.GetQuotes(appl.QuoteRequestParams{
		Recipient: appl.Recipient{
			Address: appl.Address{
				Zipcode: request.Recipient.Address.Zipcode,
			},
		},
		Volumes: vlms,
	})
	if err != nil {
		return nil, &interf.ErrorHandler{
			StatusCode: http.StatusBadRequest,
			Err:        err,
			ErrCode:    1001,
		}
	}
	carriers := make([]CreateQuoteCarrierResponse, 0)
	for _, v := range resp.Carriers {
		carriers = append(carriers, CreateQuoteCarrierResponse{
			Name:     v.Name,
			Service:  v.Service,
			Deadline: v.Deadline,
			Price:    v.Price,
		})
	}

	// success
	return &CreateQuoteResponse{
		Carrier: carriers,
	}, nil
}
