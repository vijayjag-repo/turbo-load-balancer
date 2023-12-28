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
	Serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

func (s *simpleServer) GetAddress() string {
	return s.address
}

func (s *simpleServer) IsAlive() bool {
	return true
}

func (s *simpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func createNewServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	handleError(err)

	return &simpleServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
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

func (lb *loadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *loadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address: %q\n", targetServer.GetAddress())
	targetServer.Serve(rw, r)
}

func main() {
	// Create list of servers
	servers := []Server{
		createNewServer("http://www.google.com"),
		createNewServer("http://www.duckduckgo.com"),
		createNewServer("https://www.facebook.com"),
	}

	// Create new lb
	lb := createNewLoadBalancer("8000", servers)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	})

	fmt.Println("Serving requests at http://localhost:%s\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
