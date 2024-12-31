package internals

import (
	"fmt"

	"github.com/gorilla/mux"
)

func LoadBalancer() *mux.Router {
	fmt.Println("Hello from load balancer")

	servers := Servers{
		server{Host: "localhost:8082", Port: 8082, Scheme: "http", Url: "http://localhost:8082"},
		server{Host: "localhost:3000", Port: 3000, Scheme: "http", Url: "http://localhost:3000"},
	}

	router := mux.NewRouter()

	router.HandleFunc("/{endpoint:.*}", servers.ReverseProxy)
	return router
}
