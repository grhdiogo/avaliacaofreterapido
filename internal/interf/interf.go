package interf

import (
	"net/http"

	"github.com/gorilla/mux"
)

type OutputError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorHandler error struct to show when endpoints gives an error
type ErrorHandler struct {
	// http status code
	StatusCode int
	// error
	Err error
	// internal error code
	ErrCode int
}

type AdapterFunc func(r *http.Request) (any, *ErrorHandler)

type WebServiceConfig struct {
	Prefix  string
	Version string
}

type RouteConfig struct {
	Method  string
	Path    string
	Handler AdapterFunc
}

type WebServiceImpl interface {
	GetRouter() *mux.Router
}
