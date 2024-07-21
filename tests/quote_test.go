package tests

import (
	"avaliacaofreterapido/internal/appl"
	"avaliacaofreterapido/internal/domain/quote"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestQuote(t *testing.T) {
	t.Run("create-quote", func(t *testing.T) {
		os.Setenv("FRETERAPIDO_HOST", "https://sp.freterapido.com/api/v3")
		os.Setenv("FRETERAPIDO_TOKEN", "1d52a9b6b78cf07b08586152459a5c90")
		os.Setenv("FRETERAPIDO_PLATFORM_CODE", "5AKVkHqCn")
		os.Setenv("FRETERAPIDO_DISPATHER_CEP", "29161376")
		os.Setenv("FRETERAPIDO_CNPJ", "25438296000158")
		// entrada de dados
		dataEntry := `{
			"recipient":{
				 "address":{
						"zipcode":"01311000"
				 }
			},
			"volumes":[
				 {
						"category":7,
						"amount":1,
						"unitary_weight":5,
						"price":349,
						"sku":"abc-teste-123",
						"height":0.2,
						"width":0.2,
						"length":0.2
				 },
				 {
						"category":7,
						"amount":2,
						"unitary_weight":4,
						"price":556,
						"sku":"abc-teste-527",
						"height":0.4,
						"width":0.6,
						"length":0.15
				 }
			]
	 }`

		// saída esperada
		expectedResult := `{
		 "carrier":[
				{
					 "name":"EXPRESSO FR",
					 "service":"Rodoviário",
					 "deadline":3,
					 "price":17
				},
				{
					 "name":"Correios",
					 "service":"SEDEX",
					 "deadline":1,
					 "price":20.99
				}
		 ]
	}`
		quoteService := appl.NewQuoteService()
		// transforma dados de entrada na estrutura de entrada]
		var request = appl.QuoteRequestParams{}
		err := json.Unmarshal([]byte(dataEntry), &request)
		if err != nil {
			t.Fatal(err)
		}
		// reecupera as cotas
		result, err := quoteService.GetQuotes(request)
		if err != nil {
			t.Fatal(err)
		}
		want := &quote.Entity{}
		// transforma o resultado esperado na estrutura de saída
		err = json.Unmarshal([]byte(expectedResult), want)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(result)
		if !reflect.DeepEqual(result, want) {
			t.Fatal("Resultado não é o esperado")
		}
	})
}
