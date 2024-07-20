package freterapido

import (
	"errors"
	"fmt"
)

type ReturnSimulationTypeKind int
type FreightFilterKind int
type RecipientKind int

var errMessageMapping = map[int]string{
	401: "Não autorizado (Token inválido)",                // Erro de autenticação com a API. O token informado é inválido. Ou ausência do formato de comunicação no header da requisição.
	403: "Acesso não permitido",                           // Não tem permissão para acesso ao recurso solicitado. Ou ausência de saldo para contratação de frete.
	404: "Não encontrado (Informação não localizada)",     // Recurso ou informação inválida e não pode ser encontrada.
	409: "Conflito (Atualização do status não permitido)", // A requisição não pode ser completada pois não pode ser informado um status sem seguir o fluxo de ocorrências. Os status de um frete devem seguir o fluxo normal de ocorrências.
	422: "Erro de sintaxe da requisição",                  // Houve algum problema no formato ​json​ como erro de sintaxe, valores faltando ou tipo de parâmetros inválidos.
	500: "Erro do servidor",                               // Erro interno no servidor da Frete Rápido. Você pode relatar isso para o suporte técnico da Frete Rápido através do e-mail: suporte@freterapido.com.
}

var CategoryMapping = map[int]string{
	1:   "Abrasivos",
	155: "Acessório infantil",
	109: "Acessório para decoração (com vidro)",
	110: "Acessório para decoração (sem vidro)",
	111: "Acessórios automotivos",
	69:  "Acessórios de Airsoft / Paintball",
	73:  "Acessórios de Arquearia",
	149: "Acessórios de montaria",
	70:  "Acessórios de Pesca",
	112: "Acessórios para bicicleta",
	90:  "Acessórios para celular",
	130: "Acessórios para Narguilés",
	132: "Acessórios para Tabacaria",
	2:   "Adubos / Fertilizantes",
	74:  "Alimentos não perecíveis",
	3:   "Alimentos perecíveis",
	145: "Armações de óculos",
	72:  "Arquearia",
	113: "Artesanatos (com vidro)",
	93:  "Artesanatos (sem vidro)",
	152: "Artigos de festas",
	82:  "Artigos para Camping",
	4:   "Artigos para Pesca",
	5:   "Auto Peças",
	133: "Banheira Acrílico",
	134: "Banheira de Aço Esmaltada",
	135: "Banheira Fibra de Vidro",
	156: "Banheira infantil",
	128: "Bebedouros e Purificadores",
	6:   "Bebidas / Destilados",
	114: "Bicicletas (desmontada)",
	99:  "Bijuteria",
	7:   "Brindes",
	8:   "Brinquedos",
	157: "Cadeirinha para automóvel",
	146: "Caixa d'água (Completa)",
	75:  "Caixa de embalagem",
	136: "Caixa Plástica",
	9:   "Calçados",
	115: "Cama / Mesa / Banho",
	62:  "Cargas refrigeradas/congeladas",
	158: "Carrinho de bebê",
	137: "Cartucho de Gás",
	10:  "CD / DVD / Blu-Ray",
	122: "Celulares e Smartphones",
	116: "Chapas de madeira",
	121: "Chip de celular",
	102: "Cocção Industrial",
	66:  "Colchão",
	11:  "Combustíveis / Óleos",
	12:  "Confecção",
	13:  "Cosméticos",
	14:  "Couro",
	15:  "Derivados Petróleo",
	16:  "Descartáveis",
	17:  "Editorial",
	19:  "Eletrodomésticos",
	125: "Eletrodomésticos industriais",
	150: "Eletroportáteis",
	18:  "Eletrônicos",
	20:  "Embalagens",
	138: "Equipamento oftalmológico",
	107: "Equipamentos de cozinha industrial",
	88:  "Equipamentos de Segurança / API",
	151: "Equipamentos para solda",
	84:  "Estiletes / Materiais Cortantes",
	106: "Estufa térmica",
	21:  "Explosivos / Pirotécnicos",
	126: "Expositor industrial",
	87:  "Extintores",
	23:  "Ferragens",
	24:  "Ferramentas",
	25:  "Fibras Ópticas",
	26:  "Fonográfico",
	27:  "Fotográfico",
	28:  "Fraldas / Geriátricas",
	29:  "Higiene",
	30:  "Impressos",
	31:  "Informática / Computadores",
	32:  "Instrumento Musical",
	100: "Joia",
	144: "Lentes de contato",
	86:  "Limpeza",
	77:  "Linha Branca",
	33:  "Livro(s)",
	79:  "Malas / Mochilas",
	117: "Manequins",
	104: "Maquina de algodão doce",
	105: "Maquina de chocolate",
	34:  "Materiais Escolares",
	35:  "Materiais Esportivos",
	36:  "Materiais Frágeis",
	97:  "Materiais hidráulicos / Encanamentos",
	37:  "Material de Construção",
	38:  "Material de Irrigação",
	154: "Material de Jardinagem",
	141: "Material de laboratório",
	39:  "Material Elétrico / Lâmpada(s)",
	40:  "Material Gráfico",
	41:  "Material Hospitalar",
	42:  "Material Odontológico",
	139: "Material oftalmológico",
	43:  "Material Pet Shop",
	50:  "Material Plástico",
	44:  "Material Veterinário",
	127: "Maçanetas",
	22:  "Medicamentos",
	46:  "Moto Peças",
	47:  "Mudas / Plantas",
	80:  "Máquina / Equipamentos",
	68:  "Móveis com peças de vidro",
	64:  "Móveis desmontados",
	159: "Móveis infantis",
	45:  "Móveis montados",
	129: "Narguilés",
	999: "Outros",
	48:  "Papelaria / Documentos",
	63:  "Papelão",
	49:  "Perfumaria",
	98:  "Pia / Vasos",
	83:  "Pilhas / Baterias",
	92:  "Pisos (cerâmica) / Revestimentos",
	96:  "Placa de Energia Solar",
	51:  "Pneus e Borracharia",
	95:  "Porta / Janelas (sem vidro)",
	118: "Portas / Janelas (com vidro)",
	124: "Portáteis industriais",
	85:  "Produto Químico classificado",
	52:  "Produtos Cerâmicos",
	143: "Produtos de SexShop",
	53:  "Produtos Químicos Não Classificados",
	54:  "Produtos Veterinários",
	94:  "Quadros / Molduras",
	81:  "Rações / Alimento para Animal",
	101: "Refrigeração Industrial",
	153: "Relógios",
	55:  "Revistas",
	148: "Selas e Arreios de montaria",
	56:  "Sementes",
	71:  "Simulacro de Arma / Airsoft",
	65:  "Sofá",
	57:  "Suprimentos Agrícolas / Rurais",
	131: "Tabacaria",
	147: "Tampa de Caixa d'água",
	108: "Tapeçaria / Cortinas / Persianas",
	123: "Telefonia Fixa e Sem Fio",
	142: "Tinta",
	91:  "Toldos",
	119: "Torneiras",
	67:  "Travesseiro",
	76:  "TV / Monitores",
	58:  "Têxtil",
	103: "Utensílios industriais",
	89:  "Utilidades domésticas",
	59:  "Vacinas",
	120: "Vasos de polietileno",
	60:  "Vestuário",
	61:  "Vidros / Frágil",
	78:  "Vitaminas / Suplementos nutricionais",
	140: "Óculos e acessórios",
}

// newStatusCodeErr returns an error, case not mapped, returns default
func newStatusCodeErr(sttCode int) error {
	// recover msg from map
	msg := errMessageMapping[sttCode]
	if msg != "" {
		return errors.New(msg)
	}
	// case not exists, return default
	return fmt.Errorf("Requisição falhou com código: %d", sttCode)
}

const (
	ReturnSimulationTypeFract    ReturnSimulationTypeKind = 0
	ReturnSimulationTypeCapacity ReturnSimulationTypeKind = 1
	//
	FreightFilterLowestPrice                      FreightFilterKind = 1
	FreightFilterLowestDeliveryTime               FreightFilterKind = 2
	FreightFilterLowestPriceAndLowestDeliveryTime FreightFilterKind = 3
	//
	RecipientNaturalPerson   RecipientKind = 0
	RecipientJuridicalPerson RecipientKind = 1
)

type Shipper struct {
	RegisteredNumber string `json:"registered_number"` // required
	Token            string `json:"token"`             // required
	PlatformCode     string `json:"platform_code"`     // required
}

type Recipient struct {
	Type             RecipientKind `json:"type"`                        // required
	RegisteredNumber string        `json:"registered_number,omitempty"` // optional
	StateInscription string        `json:"state_inscription,omitempty"` // optional
	Country          string        `json:"country"`                     // required
	Zipcode          int           `json:"zipcode"`                     // required
}

type DispatcherVolume struct {
	Amount        int     `json:"amount"`                   // required
	AmountVolumes int     `json:"amount_volumes,omitempty"` // optional
	Category      string  `json:"category"`                 // required
	Sku           string  `json:"sku,omitempty"`            // optional
	Tag           string  `json:"tag,omitempty"`            // optional
	Description   string  `json:"description,omitempty"`    // optional
	Height        float64 `json:"height"`                   // required
	Width         float64 `json:"width"`                    // required
	Length        float64 `json:"length"`                   // required
	UnitaryPrice  float64 `json:"unitary_price"`            // required
	UnitaryWeight float64 `json:"unitary_weight"`           // required
	Consolidate   bool    `json:"consolidate,omitempty"`    // optional
	Overlaid      bool    `json:"overlaid,omitempty"`       // optional
	Rotate        bool    `json:"rotate,omitempty"`         // optional
}

type Dispatcher struct {
	RegisteredNumber string             `json:"registered_number"`     // required
	Zipcode          int                `json:"zipcode"`               // required
	TotalPrice       float64            `json:"total_price,omitempty"` // optional
	Volumes          []DispatcherVolume `json:"volumes,omitempty"`     // required
}

type Returns struct {
	Composition  bool `json:"composition,omitempty"`   // optional
	Volumes      bool `json:"volumes,omitempty"`       // optional
	AppliedRules bool `json:"applied_rules,omitempty"` // optional
}

type CreateFreightQuotationRequest struct {
	Shipper        Shipper                    `json:"shipper"`                  // required
	Recipient      Recipient                  `json:"recipient"`                // required
	Dispatchers    []Dispatcher               `json:"dispatchers"`              // required
	Channel        string                     `json:"channel,omitempty"`        // optional
	Filter         FreightFilterKind          `json:"filter,omitempty"`         // optional
	Limit          int                        `json:"limit,omitempty"`          // optional
	Identification string                     `json:"identification,omitempty"` // optional
	Reverse        bool                       `json:"reverse,omitempty"`        // optional
	SimulationType []ReturnSimulationTypeKind `json:"simulation_type"`          // required
	Returns        *Returns                   `json:"returns,omitempty"`        // optional
}

type CreateFreightQuotationResponse struct {
	Dispatchers []struct {
		ID                         string `json:"id"`
		RequestID                  string `json:"request_id"`
		RegisteredNumberShipper    string `json:"registered_number_shipper"`
		RegisteredNumberDispatcher string `json:"registered_number_dispatcher"`
		ZipcodeOrigin              int    `json:"zipcode_origin"`
		Offers                     []struct {
			Offer          int `json:"offer"`
			SimulationType int `json:"simulation_type"`
			Carrier        struct {
				Reference        int    `json:"reference"`
				Name             string `json:"name"`
				RegisteredNumber string `json:"registered_number"`
				StateInscription string `json:"state_inscription"`
				Logo             string `json:"logo"`
			} `json:"carrier"`
			Service            string `json:"service"`
			ServiceCode        string `json:"service_code"`
			ServiceDescription string `json:"service_description"`
			DeliveryTime       struct {
				Days          int    `json:"days"`
				Hours         int    `json:"hours"`
				Minutes       int    `json:"minutes"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"delivery_time"`
			Expiration string  `json:"expiration"`
			CostPrice  float64 `json:"cost_price"`
			FinalPrice float64 `json:"final_price"`
			Weights    struct {
				Real  float64 `json:"real"`
				Cubed float64 `json:"cubed"`
				Used  float64 `json:"used"`
			} `json:"weights"`
			Composition struct {
				FreightWeight       float64 `json:"freight_weight"`
				FreightWeightExcess float64 `json:"freight_weight_excess"`
				FreightWeightVolume float64 `json:"freight_weight_volume"`
				FreightVolume       float64 `json:"freight_volume"`
				FreightMinimum      float64 `json:"freight_minimum"`
				FreightInvoice      float64 `json:"freight_invoice"`
				SubTotal1           struct {
					Daily           int `json:"daily"`
					Collect         int `json:"collect"`
					Dispatch        int `json:"dispatch"`
					Delivery        int `json:"delivery"`
					Ferry           int `json:"ferry"`
					Suframa         int `json:"suframa"`
					Tas             int `json:"tas"`
					SecCat          int `json:"sec_cat"`
					Dat             int `json:"dat"`
					AdValorem       int `json:"ad_valorem"`
					Ademe           int `json:"ademe"`
					Gris            int `json:"gris"`
					Emex            int `json:"emex"`
					Interior        int `json:"interior"`
					Capatazia       int `json:"capatazia"`
					River           int `json:"river"`
					RiverInsurance  int `json:"river_insurance"`
					Toll            int `json:"toll"`
					Other           int `json:"other"`
					OtherPerProduct int `json:"other_per_product"`
				} `json:"sub_total1"`
				SubTotal2 struct {
					Trt        int `json:"trt"`
					Tda        int `json:"tda"`
					Tde        int `json:"tde"`
					Scheduling int `json:"scheduling"`
				} `json:"sub_total2"`
				SubTotal3 struct {
					Icms int `json:"icms"`
				} `json:"sub_total3"`
			} `json:"composition"`
			OriginalDeliveryTime struct {
				Days          int    `json:"days"`
				Hours         int    `json:"hours"`
				Minutes       int    `json:"minutes"`
				EstimatedDate string `json:"estimated_date"`
			} `json:"original_delivery_time"`
			Identifier string `json:"identifier"`
		} `json:"offers"`
		Volumes []struct {
			Category      string  `json:"category"`
			Sku           string  `json:"sku"`
			Tag           string  `json:"tag"`
			Description   string  `json:"description"`
			Amount        int     `json:"amount"`
			Width         float64 `json:"width"`
			Height        float64 `json:"height"`
			Length        float64 `json:"length"`
			UnitaryWeight float64 `json:"unitary_weight"`
			UnitaryPrice  float64 `json:"unitary_price"`
			AmountVolumes float64 `json:"amount_volumes"`
			Consolidate   bool    `json:"consolidate"`
			Overlaid      bool    `json:"overlaid"`
			Rotate        bool    `json:"rotate"`
			Items         []any   `json:"items"`
		} `json:"volumes"`
	} `json:"dispatchers"`
}
