package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var nextServerIndex int32 = 0

func main() {
	var mu sync.Mutex

	// define origin server list to load balance the requests
	originServerList := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}

	loadBalancerHandler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// use mutex to prevent data race
		mu.Lock()

		// get next server to send a request to
		originServerURL, _ := url.Parse(originServerList[(nextServerIndex)%2])

		// increment next server value
		nextServerIndex++

		mu.Unlock()

		// use existing reverse proxy from httputil to route
		// a request to previously selected server url
		reverseProxy := httputil.NewSingleHostReverseProxy(originServerURL)

		reverseProxy.ServeHTTP(rw, req)
	})

	log.Fatal(http.ListenAndServe(":8080", loadBalancerHandler))
}
