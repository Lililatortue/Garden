package http

import (
	"garden/http/api"
	"net/http"
	"strings"
)

type GardenHandler struct {
	api http.Handler
	web http.Handler
}

func (g GardenHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var path = request.URL.Path
	if strings.HasPrefix(path, "/api") {
		g.api.ServeHTTP(writer, request)
	} else {
		g.web.ServeHTTP(writer, request)
	}
}

func NewGardenHandler(fspath string) *GardenHandler {
	return &GardenHandler{
		api: api.NewGardenApi(),
		web: http.FileServer(http.Dir(fspath)),
	}
}
