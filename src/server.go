package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type server interface {
	isAlive() bool
	getAddress() string
	serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

func (s *simpleServer) getAddress() string {
	return s.address
}

func (s *simpleServer) isAlive() bool {
	return true
}

func (s *simpleServer) serve(rw http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(rw, r)
}

func createNewServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	HandleError(err)

	return &simpleServer{
		address: address,
		proxy:   httputil.NewSingleHostReverseProxy(serverUrl),
	}
}
