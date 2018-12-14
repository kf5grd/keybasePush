package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	serverHost string = "127.0.0.1"
)

var instanceName string
var serverPort int

func init() {
	flag.StringVar(&instanceName, "name", "", "Set the name of this instance (required)")
	flag.IntVar(&serverPort, "port", 8617, "Set the port for the API")
	flag.Parse()
}

func main() {
	if instanceName == "" {
		fmt.Println("Error: Instance name is required.\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	router := NewRouter()

	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}
