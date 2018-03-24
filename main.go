package main

import (
	"log"
	"net"
	"net/http"

	"github.com/zserge/webview"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	go func() {
		log.Fatal(http.Serve(ln, nil))
	}()
	webview.Open("Test motherfucker", "https://google.com", 800, 600, true)
}
