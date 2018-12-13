package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	serverHost string = "127.0.0.1"
)

var instanceName string
var serverPort int

func init() {
	flag.StringVar(&instanceName, "name", "", "Set the name of this instance")
	flag.IntVar(&serverPort, "port", 8617, "Set the port for the API")
	flag.Parse()
}

func main() {
	router := NewRouter()

	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
