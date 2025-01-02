package internals

import (
	"fmt"
	"sync"

	"github.com/gorilla/mux"
)

func LoadBalancer() *mux.Router {
	fmt.Println("Hello from load balancer")

	servers := Servers{
		[]server{{Host: "localhost:8082", Port: 8082, Scheme: "http", Url: "http://localhost:8082", healthy: true},
			{Host: "localhost:3000", Port: 3000, Scheme: "http", Url: "http://localhost:3000", healthy: true}},
		&sync.Mutex{},
		0,
	}

	router := mux.NewRouter()
	servers.HealthCheck()
	router.HandleFunc("/{endpoint:.*}", servers.ReverseProxy)
	return router
}
