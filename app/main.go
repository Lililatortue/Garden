package main

import (
	"fmt"

	"garden/data/sql"
)

func main() {
	fmt.Println("Hello World")
	db, err := sql.NewDBAccess()
	fmt.Println(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	fmt.Println("Database connection established")

}
