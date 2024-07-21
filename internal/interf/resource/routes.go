package resource

import (
	"avaliacaofreterapido/internal/interf"
	"net/http"
)

var routes = []interf.RouteConfig{
	{
		Method:  http.MethodPost,
		Path:    "/quote",
		Handler: CreateQuote,
	},
}
