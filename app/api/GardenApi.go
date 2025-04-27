package api

import (
	"garden/data"
	"garden/data/sql"
	"net/http"
)

type GardenApi struct {
	*http.ServeMux
	repoManager *data.GardenService
}

func (api *GardenApi) Close() error {
	return api.repoManager.Close()
}

func NewGardenApi() *GardenApi {

	mux := &GardenApi{
		ServeMux: http.NewServeMux(),
		repoManager: data.NewRepoWith(
			sql.NewDBAccess(),
		),
	}

	mux.setRoutes()
	return mux
}
