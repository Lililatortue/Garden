package server

import (
	"garden/api"
	"net/http"
)

type GardenServer struct {
	*http.Server
}

func (server *GardenServer) Close() error {
	err := server.Handler.(*api.GardenApi).Close()
	if err != nil {
		return err
	}
	return server.Server.Close()
}

func NewGardenServer(webFsPath string, port string) *GardenServer {
	return &GardenServer{
		Server: &http.Server{
			Addr:    ":" + port,
			Handler: NewGardenHandler(webFsPath),
		},
	}
}
