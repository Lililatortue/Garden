package sql

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

func NewMockDBAccess() *DBAccess {
	db, err := sql.Open("postgres", getConnectionString("db", "db", "test_db", "5432", "garden_test"))
	if err != nil {
		panic(err)
	}
	access := &DBAccess{db}
	return access
}

func NewMockDBAccessWithSetup() *DBAccess {
	access := NewMockDBAccess()
	go access.setup()
	return access
}

func TestGetConnectionString(t *testing.T) {
	user := "test"
	passwd := "test"
	host := "test"
	port := "test"
	database := "test"
	connStr := getConnectionString(user, passwd, host, port, database)
	if connStr != "user=test password=test host=test port=test dbname=test sslmode=disable" {
		t.Errorf("connection string is not correct")
	} else {
		t.Logf("connection string is correct")
	}
}

func TestNewDBAccess(t *testing.T) {
	conn := NewDBAccess("test")
	if conn == nil {
		t.Errorf("connection is nil")
	}
	t.Logf("connection is not nil")
}

func TestDBAccess_Close(t *testing.T) {
	conn := NewMockDBAccess()
	err := conn.Close()
	if err != nil {
		t.Errorf("error closing connection")
	}
	t.Logf("connection closed")
}
