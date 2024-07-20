package freterapido

import (
	"fmt"
	"testing"
)

// TODO: Remover daqui?
const (
	cnpj    = "25438296000158"
	country = "Bra"
)

func Test_frClient_CreateFreight(t *testing.T) {
	type args struct {
		r *CreateFreightQuotationRequest
	}
	request := args{
		&CreateFreightQuotationRequest{
			Shipper: Shipper{
				RegisteredNumber: cnpj,
				Token:            "1d52a9b6b78cf07b08586152459a5c90",
				PlatformCode:     "5AKVkHqCn",
			},
			Recipient: Recipient{
				Type:    RecipientNaturalPerson,
				Country: country,
				Zipcode: 01311000,
			},
			Dispatchers: []Dispatcher{
				{
					RegisteredNumber: cnpj,
					Zipcode:          01311000,
					Volumes: []DispatcherVolume{
						{
							Amount:        1,
							Category:      CategoryMapping[7],
							UnitaryWeight: 5,
							UnitaryPrice:  349,
							Sku:           "abc-teste-123",
							Height:        0.2,
							Width:         0.2,
							Length:        0.2,
						},
					},
				},
			},
			SimulationType: []ReturnSimulationTypeKind{
				ReturnSimulationTypeCapacity,
			},
		},
	}
	response := &CreateFreightQuotationResponse{}
	//
	tests := []struct {
		name    string
		c       *frClient
		args    args
		want    *CreateFreightQuotationResponse
		wantErr bool
	}{
		{"create-freight", &frClient{
			config: Config{
				BaseUrl: "https://sp.freterapido.com/api/v3",
			},
		}, request, response, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.CreateFreight(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("frClient.CreateFreight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// resultado esperado
			// {
			// 	"dispatchers": [
			// 		{
			// 			"id": "669c140b81e7aa408fedb939",
			// 			"request_id": "669c140b81e7aa408fedb938",
			// 			"registered_number_shipper": "25438296000158",
			// 			"registered_number_dispatcher": "25438296000158",
			// 			"zipcode_origin": 365056
			// 		}
			// 	]
			// }
			// id, request_id são gerados, então não da para comparar, comparar os outros dados
			gotDispather := got.Dispatchers[0]
			requestDispather := request.r.Dispatchers[0]
			expectedResult := gotDispather.RegisteredNumberDispatcher == requestDispather.RegisteredNumber && requestDispather.Zipcode == gotDispather.ZipcodeOrigin
			if !expectedResult {
				t.Fatal("Resultado não esperado")
			}
			// passou
			fmt.Println("Resultado esperado")
		})
	}
}
