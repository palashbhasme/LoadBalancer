package internals

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func (servers *Servers) ReverseProxy(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Hello from reverse proxy")

	targetServer := roundRobin(Mu, &Servindx, servers)

	parsedURL, err := url.Parse(targetServer.Url)

	if err != nil {
		log.Println("Error parsing URL: ", err)
	}

	targetHost := parsedURL.Host
	targetScheme := parsedURL.Scheme

	r.Host = targetHost
	r.URL.Scheme = targetScheme
	r.URL.Host = targetHost
	r.RequestURI = ""

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	res, err := client.Do(r)
	if err != nil {
		log.Println("Error sending request to server: ", err)
	}

	defer func() {
		if res.Body != nil {
			res.Body.Close()
		}
	}()

	for key, values := range res.Header {
		for _, v := range values {
			w.Header().Add(key, v)
		}
	}

	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println(w, "Error wrting response", err)
	}

}

func roundRobin(mu *sync.Mutex, servindx *int, servers *Servers) server {
	mu.Lock()
	defer mu.Unlock()
	*servindx = (*servindx + 1) % len(*servers)
	targetServer := (*servers)[*servindx]
	return targetServer

}
