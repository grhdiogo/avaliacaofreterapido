package freterapido

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FrClient interface {
	// CreateFreight create a freight quotation https://dev.freterapido.com/ecommerce/cotacao_v3/
	CreateFreight(r *CreateFreightQuotationRequest) (*CreateFreightQuotationResponse, error)
}

type Config struct {
	BaseUrl string
}

type frClient struct {
	config Config
}

func (c *frClient) CreateFreight(r *CreateFreightQuotationRequest) (*CreateFreightQuotationResponse, error) {
	// parse request body
	body, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	// create request
	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/quote/simulate", c.config.BaseUrl), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("content-type", "application/json")
	// make request
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//
	decoded, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// verify response status code
	if resp.StatusCode != 200 {
		return nil, newStatusCodeErr(resp.StatusCode)
	}
	result := new(CreateFreightQuotationResponse)
	// unmarshall response
	err = json.Unmarshal(decoded, result)
	if err != nil {
		return nil, err
	}
	// success
	return result, nil
}

// NewFrClient returns a client implementation from Frete Rapido Api
func NewFrClient(config Config) FrClient {
	return &frClient{
		config: config,
	}
}
