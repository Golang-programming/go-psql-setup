package main

import (
	"fmt"
	"go-psql-setup/router"
	"log"
	"net/http"
)

func main() {
	port := 9090

	router := router.Router()

	fmt.Printf("Starting on port %d", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))

}
