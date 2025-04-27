package main

import (
	"fmt"
	"garden/server"
	"log"
)

func main() {
	fmt.Println("Hello World")
	var (
		port = "80"
		srv  = server.NewGardenServer("web", port)
	)

	fmt.Println("Starting API server on port", port)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = srv.Close()
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println("API server stopped")
	}()

}
