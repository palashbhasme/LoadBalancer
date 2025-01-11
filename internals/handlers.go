package internals

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func (servers *Servers) ReverseProxy(w http.ResponseWriter, r *http.Request) {

	proxyRequest := new(http.Request)
	*proxyRequest = *r
	proxyRequest.URL = new(url.URL)
	*proxyRequest.URL = *r.URL

	targetServer, err := servers.roundRobin()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	targetURL := fmt.Sprintf("%s://%s/%s", targetServer.Scheme, targetServer.Host, targetServer.Path)

	parsedUrl, err := url.Parse(targetURL)
	if err != nil {
		log.Println("error parsing URL: ", err)
	}

	proxyRequest.Host = parsedUrl.Host
	proxyRequest.URL.Scheme = parsedUrl.Scheme
	proxyRequest.URL.Host = parsedUrl.Host
	proxyRequest.RequestURI = ""
	proxyRequest.Header.Set("X-Forwarded-For", r.RemoteAddr)

	res, err := client.Do(proxyRequest)
	if err != nil {
		http.Error(w, "failed to connect to server", http.StatusBadGateway)
		log.Printf("error sending request to server (%s): %v", targetURL, err)
		return
	}

	defer res.Body.Close()

	for key, values := range res.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(res.StatusCode)
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Println("error writing response:", err)
	}
}

func (s *Servers) roundRobin() (server, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	unhealthy := 0

	if len(s.Servers) == 0 {
		err := errors.New("no servers available")
		return server{}, err
	}
	for {
		if unhealthy == len(s.Servers) {
			err := errors.New("all servers are down")
			return server{}, err
		}

		s.index = (s.index + 1) % len(s.Servers)
		if s.Servers[s.index].Healthy {
			return s.Servers[s.index], nil
		} else {
			unhealthy++
			continue
		}
	}

}
