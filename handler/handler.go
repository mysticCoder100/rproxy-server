package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/mysticCoder100/rproxy-server/balancer"
)

func HandleRequest(pool balancer.ServerPool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Here - Incoming request received")

		// 1. Pick the optimal backend using your Least Connections algorithm
		target := pool.NextServer()
		if target == nil {
			http.Error(w, "Bad Gateway: No backends available", http.StatusBadGateway)
			return
		}

		// 2. Safely increment the active connection count
		target.IncrementCount()
		// Ensure the count drops back down when the request completes, even if a panic occur\ws
		defer target.DecrementCount()

		fmt.Printf("Routing request -> %s (Active Conns: %d)\n", target.Url.Host, target.GetCount())

		// 3. Create the reverse proxy for the chosen backend URL
		proxy := httputil.NewSingleHostReverseProxy(target.Url)

		// 4. Forward the request and stream the response back to the client
		proxy.ServeHTTP(w, r)
	}
}
