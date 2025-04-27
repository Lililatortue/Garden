package sql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type DBAccess struct {
	*sql.DB
}

func NewDBAccess() *DBAccess {
	connStr := getConnectionString("db", "db", "localhost", "5432", "garden")
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(fmt.Errorf("error opening DB connection: %w", err))
	}
	log.Println("Connected to DB")
	access := &DBAccess{conn}
	log.Println("DB access created")
	log.Println("Setting up DB...")
	go access.setup()
	return access
}

func getConnectionString(user string, passwd string, host string, port string, database string) string {
	// NOTE: ssl should be enabled for production
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user, passwd, host, port, database)
}
