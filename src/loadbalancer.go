package main

import (
	"fmt"
	"net/http"
)

type loadBalancer struct {
	port            string
	roundRobinCount int
	servers         []server
}

func createNewLoadBalancer(port string, servers []server) *loadBalancer {
	return &loadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func (lb *loadBalancer) getNextAvailableServer() server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.isAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

func (lb *loadBalancer) serveProxy(rw http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address: %q\n", targetServer.getAddress())
	targetServer.serve(rw, r)
}
