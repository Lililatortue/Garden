package main

import (
	"fmt"
	"garden/api"
	"log"
)

func main() {
	fmt.Println("Hello World")
	var (
		port = "80"
	)

	apiSrv := api.NewApi(port)
	fmt.Println("Starting API server on port", port)
	err := apiSrv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		err = apiSrv.Close()
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Println("API server stopped")
	}()

}
