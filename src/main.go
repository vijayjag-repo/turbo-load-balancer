package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Create a list of servers
	servers := []server{
		createNewServer("https://www.espnfc.com"),
		createNewServer("https://www.duckduckgo.com"),
		createNewServer("https://www.cricinfo.com"),
	}

	// Create new lb
	lb := createNewLoadBalancer("8000", servers)

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		lb.serveProxy(rw, r)
	})

	fmt.Printf("Serving requests at http://localhost:%s \n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
