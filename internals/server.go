package internals

import (
	"net/http"
	"sync"
	"time"
)

// server struct
type server struct {
	Scheme  string
	Host    string
	Port    int
	Url     string
	healthy bool
}

type Servers struct {
	serverList []server
	mu         *sync.Mutex
	index      int
}

var client = &http.Client{
	Timeout: 30 * time.Second,
}
