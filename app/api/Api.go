package api

import (
	"net/http"
)

type Api struct {
	*http.ServeMux
	Port string
}

func NewApi(port string) *Api {
	var server *http.Server
		server = &http.Server{
			Addr: port,
		}
	return &Api{
		server: server,
		Port: port,
	}

func (api *Api) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	http.ListenAndServe(api.Port, nil)
}
