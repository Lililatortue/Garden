package api

import (
	"log"
	"net/http"
	"garden/gardentag"
	"garden/data/sql"
	"encoding/json"
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
				data: db,
			},
		},
	}

	return api
}

func (api *Api) setRoutes() {
	api.
}

func (api *ApiMux) setPushedRoute() {
	api.HandleFunc("/api/v1/push", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Println("Pushed route")
		
		var tag gardentag.GardenTag
		if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		api.data.

	})
}
