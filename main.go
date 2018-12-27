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
	if instanceName == "" {
		instanceName = KeybaseDeviceName()
	}
	router := NewRouter()

	log.Printf("Starting with instance name '%s'...\n", instanceName)
	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
