package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type (
	beckend struct {
		URL          *url.URL
		Alive        bool
		Mux          sync.RWMutex
		ReverseProxy *httputil.ReverseProxy
	}

	serverPool struct {
		Servers []beckend
		Current uint64
	}
)

func (b *beckend) isAlive() bool {
	b.Mux.RLock()
	defer b.Mux.RUnlock()
	return b.Alive
}

func (s *serverPool) getNextPeer() *beckend {
	l := len(s.Servers)
	next := (int(s.Current) + 1) % l
	for i := range l {
		idx := (next + i) % l
		if s.Servers[idx].isAlive() {
			s.Current = uint64(idx)
			return &s.Servers[idx]
		}
	}
	return nil
}

// lb load balances the incoming request
func lb(w http.ResponseWriter, r *http.Request) {
	peer := servPool.getNextPeer()
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

var servPool serverPool

func main() {
	var serverList string
	var port = 3030

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	// create http server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb),
	}

	log.Printf("Load Balancer started at :%d\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
