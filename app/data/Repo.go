package data

import (
	"fmt"

	"garden/data/sql"
)

type Repo struct {
	access *sql.DBAccess
}

func NewRepo() (*Repo, error) {
	access, err := sql.NewDBAccess()
	if err != nil {
		return nil, fmt.Errorf("Failed to create database connection: %w", err)
	}
	return &Repo{
		access: access,
	}, nil
}
