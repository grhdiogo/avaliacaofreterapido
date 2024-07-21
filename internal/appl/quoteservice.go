package appl

import (
	"avaliacaofreterapido/internal/domain/quote"
	"avaliacaofreterapido/internal/infra/cep"
	"avaliacaofreterapido/internal/infra/freterapido"
	"avaliacaofreterapido/internal/infra/postgres"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type QuoteRequestParams struct {
	Recipient Recipient
	Volumes   []Volume
}

type Recipient struct {
	Address Address
}

type Address struct {
	Zipcode string
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

type Carrier struct {
	Name     string
	Service  string
	Deadline int
	Price    float64
}

type CreateQuoteResponse struct {
	Carriers []Carrier
}

type QuoteService interface {
	GetQuotes(r QuoteRequestParams) (*CreateQuoteResponse, error)
}

type quoteServiceImpl struct {
	ctx context.Context
}

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

func (s *quoteServiceImpl) saveQuotesOnDatabase(params QuoteRequestParams, req *freterapido.CreateFreightQuotationRequest, resp *freterapido.CreateFreightQuotationResponse) error {
	// get conn
	tx, err := postgres.GetInstance().GetConn()
	if err != nil {
		return errors.New("Falhar ao conectar com banco de dados")
	}
	//
	rep := postgres.NewQuoteRepository(s.ctx, tx)
	rawResponse, err := json.Marshal(resp)
	if err != nil {
		return errors.New("Falha ao transformar dados de resposta da cotação")
	}
	//
	volumes := make([]quote.Volume, 0)
	for _, v := range params.Volumes {
		//
		volumes = append(volumes, quote.Volume{
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
	rawReq, _ := json.Marshal(req)
	// store
	err = rep.Store(quote.Entity{
		ID:      uuid.New().String(),
		CpfCnpj: req.Shipper.RegisteredNumber,
		Address: quote.Address{
			Cep: params.Recipient.Address.Zipcode,
		},
		RawResponse: rawResponse,
		RawRequest:  rawReq,
		Volumes:     volumes,
	})
	if err != nil {
		return errors.New("Falha ao salvar cotação")
	}
	// commit transaction
	err = tx.Commit(s.ctx)
	if err != nil {
		return errors.New("Falha ao salvar dados no banco")
	}
	// sucess
	return nil
}

func (s *quoteServiceImpl) GetQuotes(params QuoteRequestParams) (*CreateQuoteResponse, error) {
	// valida os parametros de entrada
	err := s.validate(params)
	if err != nil {
		return nil, err
	}
	// rep
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
	// reqeust to freterapido api
	req := &freterapido.CreateFreightQuotationRequest{
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
	}
	//
	response, err := frclient.CreateFreight(req)
	if err != nil {
		return nil, err
	}
	// save
	err = s.saveQuotesOnDatabase(params, req, response)
	if err != nil {
		return nil, err
	}
	//
	result := new(CreateQuoteResponse)
	// verify if have offers
	if len(response.Dispatchers) > 0 {
		dispather := response.Dispatchers[0]
		for _, v := range dispather.Offers {
			result.Carriers = append(result.Carriers, Carrier{
				Name:     v.Carrier.Name,
				Service:  v.Service,
				Deadline: v.DeliveryTime.Days,
				Price:    v.FinalPrice,
			})
		}
	}
	return result, nil
}

func NewQuoteService(ctx context.Context) QuoteService {
	return &quoteServiceImpl{
		ctx: ctx,
	}
}
