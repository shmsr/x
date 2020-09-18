package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net"
	"net/http"
)

const (
	flagHost = "host"
	flagPort = "port"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = "8080"
)

const (
	usageHost = "enter host"
	usagePort = "enter port"
)

var (
	host string
	port string
)

func serve(u string, m *http.ServeMux) {
	log.Printf("Listening on %s\n", u)
	log.Fatalln(http.ListenAndServe(u, m))
}

func main() {
	flag.StringVar(&host, flagHost, defaultHost, usageHost)
	flag.StringVar(&port, flagPort, defaultPort, usagePort)
	flag.Parse()

	url := net.JoinHostPort(host, port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Served from: %s\n", url)
		fmt.Fprintf(w, "Server says: Hello, %q!", html.EscapeString(r.RemoteAddr))
	})

	serve(url, mux)
}
