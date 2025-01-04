package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/palashbhasme/loadbalancer/internals"
	"github.com/palashbhasme/loadbalancer/utils"
)

func main() {

	servers, err := utils.LoadConfig()
	if err != nil {
		fmt.Println("error loading config: ", err)
		os.Exit(1)
	}

	port := 8000
	address := fmt.Sprintf(":%d", port)
	fmt.Printf("Load Balancer listening on %s\n", address)

	router := internals.LoadBalancer(servers)

	err = http.ListenAndServe(address, router)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		os.Exit(1)
	}

}
