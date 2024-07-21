package freterapido

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_frClient_CreateFreight(t *testing.T) {
	type args struct {
		r *CreateFreightQuotationRequest
	}
	request := args{
		&CreateFreightQuotationRequest{
			Shipper: Shipper{
				RegisteredNumber: "25438296000158",
				Token:            "1d52a9b6b78cf07b08586152459a5c90",
				PlatformCode:     "5AKVkHqCn",
			},
			Recipient: Recipient{
				Type:    RecipientNaturalPerson,
				Country: "BRA",
				Zipcode: 01311000,
			},
			Dispatchers: []Dispatcher{
				{
					RegisteredNumber: "25438296000158",
					Zipcode:          29161376,
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
				ReturnSimulationTypeFract,
			},
		},
	}
	response := &CreateFreightQuotationResponse{
		Dispatchers: []CreateFreightQuotationResponseDispatcher{
			{
				Offers: []CreateFreightQuotationResponseDispatcherOffer{
					{
						Offer:          1,
						SimulationType: 0,
						Carrier: CreateFreightQuotationResponseDispatcherOfferCarrier{
							Reference:        281,
							Name:             "CORREIOS",
							RegisteredNumber: "34028316000103",
							StateInscription: "ISENTO",
							Logo:             "https://s3.amazonaws.com/public.prod.freterapido.uploads/transportadora/foto-perfil/34028316000103.png",
						},
						Service: "PAC",
						DeliveryTime: CreateFreightQuotationResponseDispatcherOfferDeliveryTime{
							Days:          5,
							EstimatedDate: "2024-07-26",
						},
						Expiration: "2024-07-26T00:00:00Z",
						CostPrice:  29.44,
						FinalPrice: 29.44,
						Weights: CreateFreightQuotationResponseDispatcherOfferWeights{
							Real: 5,
						},
						OriginalDeliveryTime: CreateFreightQuotationResponseDispatcherOfferOriginalDeliveryTime{
							Days:          5,
							EstimatedDate: "2024-07-26",
						},
						Identifier:   "03298",
						HomeDelivery: true,
						Modal:        "Rodoviário",
						Esg: CreateFreightQuotationResponseDispatcherOfferEsg{
							Co2EmissionEstimate: 294.292254424126,
						},
					},
				},
			},
		},
	}
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
			if len(got.Dispatchers) != 1 && reflect.DeepEqual(got.Dispatchers[0].Offers, tt.want.Dispatchers[0].Offers) {
				t.Fatal("Resultado não esperado")
			}
			// passou
			fmt.Println("Resultado esperado")
		})
	}
}
