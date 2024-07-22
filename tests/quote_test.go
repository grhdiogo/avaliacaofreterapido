package tests

import (
	"avaliacaofreterapido/internal/appl"
	"avaliacaofreterapido/internal/infra/postgres"
	"avaliacaofreterapido/internal/interf/resource"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func init() {
	os.Setenv("FRETERAPIDO_HOST", "https://sp.freterapido.com/api/v3")
	os.Setenv("FRETERAPIDO_TOKEN", "1d52a9b6b78cf07b08586152459a5c90")
	os.Setenv("FRETERAPIDO_PLATFORM_CODE", "5AKVkHqCn")
	os.Setenv("FRETERAPIDO_DISPATHER_CEP", "29161376")
	os.Setenv("FRETERAPIDO_CNPJ", "25438296000158")
	//
	os.Setenv("POSTGRES_USER", "postgres")
	os.Setenv("DATABASE_PORT", "5432")
	os.Setenv("POSTGRES_DB", "dbfr")
	os.Setenv("POSTGRES_PASSWORD", "root")
	os.Setenv("PORT", "8000")

	pgconfig := postgres.Config{
		Host: os.Getenv("DATABASE_HOST"),
		User: os.Getenv("POSTGRES_USER"),
		Port: os.Getenv("DATABASE_PORT"),
		DBNm: os.Getenv("POSTGRES_DB"),
		Pswd: os.Getenv("POSTGRES_PASSWORD"),
	}
	postgres.SetConfiguration(pgconfig)
	inst := postgres.GetInstance()
	err := inst.Init(false)
	if err != nil {
		panic(err)
	}
}

func TestQuote(t *testing.T) {
	t.Run("create-quote-sucess", func(t *testing.T) {
		//
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

		//
		expectedResult := `{
		"carrier": [
			{
				"name": "JADLOG",
				"service": ".PACKAGE",
				"deadline": 3,
				"price": 31.82
			},
			{
				"name": "AZUL CARGO",
				"service": "Convencional",
				"deadline": 2,
				"price": 41.82
			},
			{
				"name": "PRESSA FR (TESTE)",
				"service": "Normal",
				"deadline": 0,
				"price": 58.95
			},
			{
				"name": "BTU BRASPRESS",
				"service": "Normal",
				"deadline": 5,
				"price": 93.35
			}
		]
	}`
		//
		buf := bytes.NewBuffer([]byte(dataEntry))
		//
		result, errorHandler := resource.CreateQuote(&http.Request{
			Method: "POST",
			Body:   io.NopCloser(buf),
		})
		if errorHandler != nil {
			t.Fatal(errorHandler.Err)
		}
		result1 := new(resource.CreateQuoteResponse)
		result2 := result.(*resource.CreateQuoteResponse)
		err := json.Unmarshal([]byte(expectedResult), result1)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(result1, result2) {
			t.Fatal("Resultado não é o esperado")
		}
	})
}

func TestQuoteFail(t *testing.T) {
	t.Run("create-quote-error-quantity", func(t *testing.T) {
		//
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
						"height":0,
						"width":0,
						"length":0
				 }
			]
	 }`

		expectedError := "Altura do 1º volume é inválido; Largura do 1º volume é inválido; Tamanho do 1º volume é inválido"

		//
		//
		buf := bytes.NewBuffer([]byte(dataEntry))
		//
		_, errorHandler := resource.CreateQuote(&http.Request{
			Method: "POST",
			Body:   io.NopCloser(buf),
		})
		if errorHandler != nil && errorHandler.Err.Error() != expectedError {
			t.Fatal("Resultado não esperado")
		}

	})
}

func TestMetricsSuccess(t *testing.T) {
	t.Run("get-metrics-success", func(t *testing.T) {
		// expected result
		expected := &appl.Metric{
			CarriersMetrics: map[string]appl.CarrierMetrict{
				"JADLOG": {
					ResultQuantity:       2,
					TotalValue:           33.1,
					MostExpensiveFreight: 18.48,
					CheaperFreight:       14.62,
				},
				"AZUL CARGO": {
					ResultQuantity:       2,
					TotalValue:           48.72,
					MostExpensiveFreight: 27.72,
					CheaperFreight:       21,
				},
				"CORREIOS": {
					ResultQuantity:       2,
					TotalValue:           94.66,
					MostExpensiveFreight: 65.22,
					CheaperFreight:       29.44,
				},
				"PRESSA FR (TESTE)": {
					ResultQuantity:       2,
					TotalValue:           114.97,
					MostExpensiveFreight: 57.82,
					CheaperFreight:       57.15,
				},
				"BTU BRASPRESS": {
					ResultQuantity:       2,
					TotalValue:           174.15,
					MostExpensiveFreight: 93.23,
					CheaperFreight:       80.92,
				},
			},
		}
		// make two quotations
		app := appl.NewQuoteService(context.Background())
		//
		_, err := app.GetQuotes(appl.QuoteRequestParams{
			Recipient: appl.Recipient{
				Address: appl.Address{
					Zipcode: "01311000",
				},
			},
			Volumes: []appl.Volume{
				{
					Category:      7,
					Amount:        1,
					UnitaryWeight: 5,
					Price:         349,
					Sku:           "abc-teste-123",
					Height:        0.2,
					Width:         0.2,
					Length:        0.2,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		_, err = app.GetQuotes(appl.QuoteRequestParams{
			Recipient: appl.Recipient{
				Address: appl.Address{
					Zipcode: "01311000",
				},
			},
			Volumes: []appl.Volume{
				{
					Category:      7,
					Amount:        2,
					UnitaryWeight: 4,
					Price:         556,
					Sku:           "abc-teste-527",
					Height:        0.4,
					Width:         0.6,
					Length:        0.15,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		// recover metrics from last two quotations
		metrics, err := app.Metricts(2)
		if err != nil {
			t.Fatal(err)
		}
		// check if metrics match the expected
		if !(reflect.DeepEqual(metrics, expected)) {
			t.Fatal("Resultado não essperado")
		}
	})
}

func TestMetricsFail(t *testing.T) {
	t.Run("get-metrics-fail", func(t *testing.T) {
		// expected result
		expected := &appl.Metric{
			CarriersMetrics: map[string]appl.CarrierMetrict{
				"JADLOG": {
					ResultQuantity:       0,
					TotalValue:           0,
					MostExpensiveFreight: 0,
					CheaperFreight:       0,
				},
				"AZUL CARGO": {
					ResultQuantity:       0,
					TotalValue:           0,
					MostExpensiveFreight: 0,
					CheaperFreight:       0,
				},
				"CORREIOS": {
					ResultQuantity:       0,
					TotalValue:           0,
					MostExpensiveFreight: 0,
					CheaperFreight:       0,
				},
				"PRESSA FR (TESTE)": {
					ResultQuantity:       0,
					TotalValue:           0,
					MostExpensiveFreight: 0,
					CheaperFreight:       0,
				},
				"BTU BRASPRESS": {
					ResultQuantity:       0,
					TotalValue:           0,
					MostExpensiveFreight: 0,
					CheaperFreight:       0,
				},
			},
		}
		// make only one quotations
		app := appl.NewQuoteService(context.Background())
		//
		_, err := app.GetQuotes(appl.QuoteRequestParams{
			Recipient: appl.Recipient{
				Address: appl.Address{
					Zipcode: "01311000",
				},
			},
			Volumes: []appl.Volume{
				{
					Category:      7,
					Amount:        1,
					UnitaryWeight: 5,
					Price:         349,
					Sku:           "abc-teste-123",
					Height:        0.2,
					Width:         0.2,
					Length:        0.2,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		// recover metrics from last two quotations
		metrics, err := app.Metricts(2)
		if err != nil {
			t.Fatal(err)
		}
		// pass if metrics diff
		if reflect.DeepEqual(metrics, expected) {
			t.Fatal("Resultado não essperado")
		}
	})
}
