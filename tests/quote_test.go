package tests

import (
	"avaliacaofreterapido/internal/infra/postgres"
	"avaliacaofreterapido/internal/interf/resource"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestQuote(t *testing.T) {
	t.Run("create-quote1", func(t *testing.T) {
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
				 }
			]
	 }`

		//
		expectedResult := `{
		 "carrier":[
				{
					 "name":"CORREIOS",
					 "service":"PAC",
					 "deadline":5,
					 "price":29.44
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
		if err != nil {
			t.Fatal(errorHandler.Err)
		}
		result1 := new(resource.CreateQuoteResponse)
		result2 := result.(*resource.CreateQuoteResponse)
		err = json.Unmarshal([]byte(expectedResult), result1)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(result1)
		fmt.Println(result2)
		if !reflect.DeepEqual(result1, result2) {
			t.Fatal("Resultado não é o esperado")
		}
		fmt.Println("Resultado esperado")
	})
}
