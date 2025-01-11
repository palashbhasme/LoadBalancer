package internals

import (
	"net/http"
	"sync"
	"time"
)

// server struct
type server struct {
	Scheme  string `yaml:"Scheme"`
	Host    string `yaml:"Host"`
	Port    int    `yaml:"Port"`
	Path    string `yaml:"Path"`
	Healthy bool   `yaml:"Healthy"`
}

type Servers struct {
	Servers []server `yaml:"servers"`
	mu      *sync.Mutex
	index   int
}

func (servers *Servers) SetIndex(index int) {
	servers.index = index
}
func (servers *Servers) SetMu(mu *sync.Mutex) {
	servers.mu = mu
}

var client = &http.Client{
	Timeout: 30 * time.Second,
}
