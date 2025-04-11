package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBAccess struct {
	*sql.DB
}

func NewDBAccess(db *sql.DB) (*DBAccess, error) {
	connStr := getConnectionString("user", "pass", "addr", "service") // TODO: write corect connection string
	conn, err := sql.Open("postgress", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening DB connection: %w", err)
	}
	access := &DBAccess{conn}
	go access.setup()
	return access, nil
}

func getConnectionString(user string, passwd string, adrr string, service string) string {
	// NOTE: ssl should be enabled for production
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		user, passwd, adrr, service)
}
