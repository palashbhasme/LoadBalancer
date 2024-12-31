package internals

import (
	"sync"
)

// server struct
type server struct {
	Scheme string
	Host   string
	Port   int
	Url    string
}

// list of servers
type Servers []server

var (
	Servindx             = 0
	Mu       *sync.Mutex = &sync.Mutex{}
)
