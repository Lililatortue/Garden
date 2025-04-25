package data

import (
	"fmt"
	"garden/data/sql"
)

type Repo struct {
	Access *sql.DBAccess
}

func NewRepo() (*Repo, error) {
	access, err := sql.NewDBAccess()
	if err != nil {
		return nil, fmt.Errorf("Failed to create database connection: %w", err)
	}
	return &Repo{
		Access: access,
	}, nil
}

func NewRepoWith(access *sql.DBAccess) *Repo {
	return &Repo{
		Access: access,
	}
}

func (repo *Repo) Close() {
	repo.Access.Close()
}
