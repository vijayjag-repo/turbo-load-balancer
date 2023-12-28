package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

type Server interface {
	IsAlive() bool
	GetAddress() string
	Serve(rw *http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

func createNewServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	handleError(err)

	return &simpleServer{
		addr:  address,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type loadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func createNewLoadBalancer(port string, servers []Server) *loadBalancer {
	return &loadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
