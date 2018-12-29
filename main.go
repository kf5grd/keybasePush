package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	instanceName = strings.ToLower(instanceName)
	instanceName = strings.TrimSpace(instanceName)
	CreateMissingChannels(instanceName)

	router := NewRouter()

	log.Printf("Starting with instance name '%s'\n", instanceName)
	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

func CreateMissingChannels(instanceName string) {
	// need queue and input channels
	neededChannels := []string{
		fmt.Sprintf("__%s_queue", instanceName),
		fmt.Sprintf("__%s_input", instanceName),
	}

	// get existing dev channels
	existingChannels := make(map[string]string)
	for _, c := range GetDevChannels() {
		existingChannels[c] = ""
	}

	// create any dev channel that's needed and doesn't exist yet
	for _, devChan := range neededChannels {
		if _, ok := existingChannels[devChan]; !ok {
			log.Printf("Creating missing dev channel '%s'\n", devChan)
			if err := CreateDevChannel(KeybaseUsername(), devChan); err != nil {
				panic(err)
			}
		}
	}
}
