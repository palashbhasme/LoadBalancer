package main

import (
	"fmt"
	"net/http"

	"www.github.com/palashbhasme/loadbalancer/internals"
)

func main() {

	port := 8000
	address := fmt.Sprintf(":%d", port)
	fmt.Printf("Load Balancer listning on %s\n", address)

	router := internals.LoadBalancer()

	err := http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	}

}
