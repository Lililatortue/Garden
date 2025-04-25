package api

import (
	"encoding/json"
	"garden/data/sql"
	"garden/types"
	"log"
	"net/http"
)

type Api struct {
	*http.Server
}

type ApiMux struct {
	*http.ServeMux
	data *sql.DBAccess
}

func NewApi(port string) *Api {
	db, err := sql.NewDBAccess()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	api := &Api{
		Server: &http.Server{
			Addr: ":" + port,
			Handler: ApiMux{
				ServeMux: http.NewServeMux(),
				data:     db,
			},
		},
	}

	return api
}

func (api *ApiMux) setRoutes() {
	api.setPushedRoute()
}

func (api *ApiMux) setPushedRoute() {
	api.HandleFunc("/api/v1/push", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Println("Pushed route")

		var tag types.GardenTag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		//api.data

	})
}
