package main

import (
	"fmt"
	"log"
	"net"
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

func isAlive(b beckend) bool {
	if !b.Alive {
		return false
	}
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", b.URL.Host, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	b.Alive = false
	return true
}

func (sp *serverPool) add(serv string) {
	servUrl, _ := url.Parse(serv)

	sp.Servers = append(sp.Servers, beckend{
		URL:   servUrl,
		Alive: true,
	})

	fmt.Printf("added %v\n", serv)
}

func (s *serverPool) getNextPeer() *beckend {
	l := len(s.Servers)
	next := (int(s.Current) + 1) % l
	for i := range l {
		idx := (next + i) % l
		if isAlive(s.Servers[idx]) {
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
		newPeer.Alive = false
		newPeer = serverList.getNextPeer()
		originServerURL = newPeer.URL
	}

	// use existing reverse proxy from httputil to route
	// a request to previously selected server url
	reverseProxy := httputil.NewSingleHostReverseProxy(originServerURL)

	reverseProxy.ServeHTTP(rw, req)
})

func Adder(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("TODO"))
		if r.Header.Get("TODO") == "Add me" {
			serverList.add(r.Header.Get("serv"))

			fmt.Println(serverList)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	var port = 8080
	var mux = Adder(loadBalancerHandler)

	// create http server
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}

	log.Printf("Load Balancer started at :%v\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
