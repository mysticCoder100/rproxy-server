package main

import (
	"fmt"
	"net/http"

	"github.com/mysticCoder100/rproxy-server/balancer"
	"github.com/mysticCoder100/rproxy-server/handler"
)

func main() {
	pool, error := balancer.InitServers()

	if error != nil {
		panic(error)
	}

	http.HandleFunc("/", handler.HandleRequest(pool))

	fmt.Println("Reverse proxy running on :5000...")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
