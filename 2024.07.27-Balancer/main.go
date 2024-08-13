package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type (
	beckend struct {
		URL   *url.URL
		Alive bool
	}

	serverPool struct {
		Servers []beckend
		Current uint32
	}
)

func (sp *serverPool) add(serv string) {
	servUrl, _ := url.Parse(serv)

	sp.Servers = append(sp.Servers, beckend{
		URL:   servUrl,
		Alive: true,
	})
}

func (s *serverPool) getNextPeer() *beckend {
	l := len(s.Servers)
	next := (int(s.Current) + 1) % l
	for i := range l {
		idx := (next + i) % l
		if s.Servers[idx].Alive {
			s.Current = uint32(idx)
			return &s.Servers[idx]
		}
	}
	return nil
}

var serverList serverPool

var loadBalancerHandler = http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

	// get next server to send a request to
	newPeer := serverList.getNextPeer()
	originServerURL := newPeer.URL

	if originServerURL == nil {
		log.Fatal("no url")
	}

	// use existing reverse proxy from httputil to route
	// a request to previously selected server url
	reverseProxy := httputil.NewSingleHostReverseProxy(originServerURL)

	reverseProxy.ServeHTTP(rw, req)
})

func main() {
	var port = 8080
	serverList.Current = 0

	serverList.add("http://localhost:8081")
	serverList.add("http://localhost:8082")

	if len(serverList.Servers) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	// create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      loadBalancerHandler,
	}

	log.Printf("Load Balancer started at :%v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
