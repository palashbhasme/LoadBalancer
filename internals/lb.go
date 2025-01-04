package internals

import (
	"fmt"

	"github.com/gorilla/mux"
)

func LoadBalancer(servers *Servers) *mux.Router {
	fmt.Println("Hello from load balancer")

	router := mux.NewRouter()
	servers.HealthCheck()
	router.HandleFunc("/{endpoint:.*}", servers.ReverseProxy)
	return router
}
