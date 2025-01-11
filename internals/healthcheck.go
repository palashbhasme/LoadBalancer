package internals

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func (s *Servers) HealthCheck() {

	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {

			s.doChecks()

		}
	}()
}

func (s *Servers) doChecks() {

	var wg sync.WaitGroup

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.Servers {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			targetServer := s.Servers[index]
			targetURL := fmt.Sprintf("%s://%s/%s", targetServer.Scheme, targetServer.Host, targetServer.Path)

			resp, err := client.Get(targetURL)
			if err != nil {
				s.Servers[index].Healthy = false
				log.Println("server is unresponsive: ", targetURL)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				s.Servers[index].Healthy = true
				log.Printf("server %s is healthy", targetURL)
			} else {
				s.Servers[index].Healthy = false
				log.Printf("server %s returned status code %d", targetURL, resp.StatusCode)
			}
		}(i)

	}
	wg.Wait()
}
