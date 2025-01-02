package internals

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func (servers *Servers) ReverseProxy(w http.ResponseWriter, r *http.Request) {

	proxyRequest := new(http.Request)
	*proxyRequest = *r
	proxyRequest.URL = new(url.URL)
	*proxyRequest.URL = *r.URL

	targetURL, err := servers.roundRobin()
	if err != nil {
		http.Error(w, "all servers are down", http.StatusInternalServerError)
	}

	parsedUrl, err := url.Parse(targetURL.Url)
	if err != nil {
		log.Println("error parsing URL: ", err)
	}

	proxyRequest.Host = parsedUrl.Host
	proxyRequest.URL.Scheme = parsedUrl.Scheme
	proxyRequest.URL.Host = parsedUrl.Host
	proxyRequest.RequestURI = ""

	res, err := client.Do(proxyRequest)
	if err != nil {
		log.Println("error sending request to server: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	for key, values := range res.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println(w, "error wrting response", err)
	}

}
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

	for i := range s.serverList {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			endpoint := fmt.Sprint(s.serverList[index].Url, "/test")

			resp, err := client.Get(endpoint)
			if err != nil {
				s.serverList[index].healthy = false
				log.Println("server is unresponsive: ", s.serverList[index].Url)
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				s.serverList[index].healthy = true
				log.Printf("server %s is healthy", s.serverList[index].Url)
			} else {
				s.serverList[index].healthy = false
				log.Printf("server %s returned status code %d", s.serverList[index].Url, resp.StatusCode)
			}
		}(i)

	}
	wg.Wait()
}

func (s *Servers) roundRobin() (server, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	unhealthy := 0

	for {
		if unhealthy == len(s.serverList) {
			log.Println("All servers are down")
			err := errors.New("all servers are down")
			return server{}, err
		}

		s.index = (s.index + 1) % len(s.serverList)
		if s.serverList[s.index].healthy {
			return s.serverList[s.index], nil
		} else {
			unhealthy++
			continue
		}
	}

}
