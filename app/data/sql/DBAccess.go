package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DBAccess struct {
	*sql.DB
}

func NewDBAccess() (*DBAccess, error) {
	connStr := getConnectionString("db", "db", "localhost", "5432", "garden")
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening DB connection: %w", err)
	}
	fmt.Println("Connected to DB")
	access := &DBAccess{conn}
	fmt.Println("DB access created")
	fmt.Println("Setting up DB...")
	go access.setup()
	return access, nil
}

func getConnectionString(user string, passwd string, host string, port string, database string) string {
	// NOTE: ssl should be enabled for production
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user, passwd, host, port, database)
}
