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
	CreateMissingChannels()
	if instanceName == "" {
		instanceName = KeybaseDeviceName()
	}
	instanceName = strings.ToLower(instanceName)
	router := NewRouter()

	log.Printf("Starting with instance name '%s'...\n", instanceName)
	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

func CreateMissingChannels() {
	deviceName := KeybaseDeviceName()
	neededChannels := []string{fmt.Sprintf("__%s_queue", deviceName), fmt.Sprintf("__%s_input", deviceName)}
	existingChannels := make(map[string]string)
	for _, c := range GetDevChannels() {
		existingChannels[c] = ""
	}

	for _, devChan := range neededChannels {
		if _, ok := existingChannels[devChan]; !ok {
			fmt.Printf("Creating missing dev channel: %s... ", devChan)
			if err := CreateDevChannel(KeybaseUsername(), devChan); err != nil {
				fmt.Printf("Error: %s\n", err)
			} else {
				fmt.Printf("Success!\n")
			}
		}
	}
}
