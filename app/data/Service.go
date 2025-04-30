package data

import (
	"garden/data/sql"
)

type GardenService struct {
	Access *sql.DBAccess
}

func NewGardenService() *GardenService {
	return &GardenService{
		Access: sql.NewDBAccess(),
	}
}

func NewRepoWith(access *sql.DBAccess) *GardenService {
	return &GardenService{
		Access: access,
	}
}

func (gs *GardenService) Close() error {
	return gs.Access.Close()
}
