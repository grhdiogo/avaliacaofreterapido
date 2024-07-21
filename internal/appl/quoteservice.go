package appl

import (
	"avaliacaofreterapido/internal/domain/quote"
	"avaliacaofreterapido/internal/infra/cep"
	"avaliacaofreterapido/internal/infra/freterapido"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	Price         float64 `json:"price"`
	Sku           string  `json:"sku"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
}

type QuoteService interface {
	GetQuotes(r QuoteRequestParams) (*quote.Entity, error)
}

type quoteServiceImpl struct{}

func (s *quoteServiceImpl) validate(p QuoteRequestParams) error {
	var errs = make([]string, 0)
	if p.Recipient.Address.Zipcode == "" || !cep.CheckZipCode(p.Recipient.Address.Zipcode) {
		errs = append(errs, "CEP inválido")
	}
	if len(p.Volumes) == 0 {
		errs = append(errs, "Ao menos 1(um) volume deve ser passado")
	}
	for index, v := range p.Volumes {
		i := index + 1
		//
		if freterapido.CategoryMapping[v.Category] == "" {
			errs = append(errs, fmt.Sprintf("Categoria do %dº volume é inválido", i))
		}
		// Amount
		if v.Amount <= 0 {
			errs = append(errs, fmt.Sprintf("Quantidade do %dº volume é inválido", i))
		}
		// UnitaryWeight
		if v.UnitaryWeight <= 0 {
			errs = append(errs, fmt.Sprintf("Peso unitário do %dº volume é inválido", i))
		}
		// Price
		if v.Price <= 0 {
			errs = append(errs, fmt.Sprintf("Preço do %dº volume é inválido", i))
		}
		// Sku
		if len(v.Sku) > 255 {
			errs = append(errs, fmt.Sprintf("Quantidade de caracteres de Sku do %dº volume é muito grande", i))
		}
		// Height
		if v.Height <= 0 {
			errs = append(errs, fmt.Sprintf("Altura do %dº volume é inválido", i))
		}
		// Width
		if v.Width <= 0 {
			errs = append(errs, fmt.Sprintf("Largura do %dº volume é inválido", i))
		}
		// Length
		if v.Length <= 0 {
			errs = append(errs, fmt.Sprintf("Tamanho do %dº volume é inválido", i))
		}
	}
	// case exist error
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	// success
	return nil
}

func (s *quoteServiceImpl) GetQuotes(params QuoteRequestParams) (*quote.Entity, error) {
	// valida os parametros de entrada
	err := s.validate(params)
	if err != nil {
		return nil, err
	}
	// instanciar repositório
	frclient := freterapido.NewFrClient(freterapido.Config{
		BaseUrl: os.Getenv("FRETERAPIDO_HOST"),
	})
	volumes := make([]freterapido.DispatcherVolume, 0)
	for _, v := range params.Volumes {
		volumes = append(volumes, freterapido.DispatcherVolume{
			Amount:        v.Amount,
			Category:      freterapido.CategoryMapping[v.Category],
			UnitaryWeight: float64(v.UnitaryWeight),
			UnitaryPrice:  v.Price,
			Sku:           v.Sku,
			Height:        v.Height,
			Width:         v.Width,
			Length:        v.Length,
		})
	}
	zipcodeRecipient, _ := strconv.Atoi(params.Recipient.Address.Zipcode)
	zipcodeDispather, _ := strconv.Atoi(os.Getenv("FRETERAPIDO_DISPATHER_CEP"))
	//
	response, err := frclient.CreateFreight(&freterapido.CreateFreightQuotationRequest{
		Shipper: freterapido.Shipper{
			RegisteredNumber: os.Getenv("FRETERAPIDO_CNPJ"),
			Token:            os.Getenv("FRETERAPIDO_TOKEN"),
			PlatformCode:     os.Getenv("FRETERAPIDO_PLATFORM_CODE"),
		},
		Recipient: freterapido.Recipient{
			Type:    freterapido.RecipientNaturalPerson,
			Country: "BRA",
			Zipcode: zipcodeRecipient,
		},
		Dispatchers: []freterapido.Dispatcher{
			{
				RegisteredNumber: os.Getenv("FRETERAPIDO_CNPJ"),
				Zipcode:          zipcodeDispather,
				Volumes:          volumes,
			},
		},
		SimulationType: []freterapido.ReturnSimulationTypeKind{
			freterapido.ReturnSimulationTypeFract,
		},
	})
	if err != nil {
		return nil, err
	}
	result := new(quote.Entity)
	// verifica se tem ofertas
	if len(response.Dispatchers) > 0 {
		dispather := response.Dispatchers[0]
		for _, v := range dispather.Offers {
			result.Carrier = append(result.Carrier, quote.Carrier{
				Name:     v.Carrier.Name,
				Service:  v.Service,
				Deadline: v.DeliveryTime.Days,
				Price:    v.FinalPrice,
			})
		}
	}
	return result, nil
}

func NewQuoteService() QuoteService {
	return &quoteServiceImpl{}
}
