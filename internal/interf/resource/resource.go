package resource

import (
	"avaliacaofreterapido/internal/interf"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type webApplication struct {
}

func createFunc(cfg interf.RouteConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := cfg.Handler(r)
		if err != nil {
			w.WriteHeader(err.StatusCode)
			//
			j, _ := json.Marshal(&interf.OutputError{
				Code:    err.ErrCode,
				Message: err.Err.Error(),
			})
			// write error
			w.Write(j)
			return
		}
		// case success
		w.WriteHeader(200)
		// recover success struct, case error, return empty
		j, _ := json.Marshal(result)
		w.Write(j)
	}
}

func (w webApplication) GetRouter() *mux.Router {
	r := mux.NewRouter()
	// prefix
	// subrouter
	router := r.NewRoute().Subrouter()
	for i := 0; i < len(routes); i++ {
		router.HandleFunc(routes[i].Path, createFunc(routes[i])).Methods(routes[i].Method)
	}
	return router
}

func NewWebApplication(cfg interf.WebServiceConfig) interf.WebServiceImpl {
	return &webApplication{}
}
