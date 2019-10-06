package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"samhofi.us/x/keybase"
)

const (
	// IP address to serve rest API on
	serverHost string = "127.0.0.1"
)

// Initialize variables for flags
var instanceName string
var serverPort int

// Keybase API
var k = keybase.NewKeybase()

// Event Commands
var eventCommands = make(map[string]eventCommand)

var wg sync.WaitGroup

func init() {
	// Parse flags
	flag.StringVar(&instanceName, "name", "", "Set the name of this instance")
	flag.IntVar(&serverPort, "port", 8617, "Set the port for the API")
	flag.Parse()
}

func apiService() {
	// Create a new router and start the rest API service
	log.Printf("Starting web service")
	router := newRouter()
	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

func main() {
	// If -name flag isn't set, use device name from keybase
	if instanceName == "" {
		instanceName = k.Device
	}
	// Instance name forced to lowercase with no leading or trailing spaces
	instanceName = strings.ToLower(instanceName)
	instanceName = strings.TrimSpace(instanceName)

	log.Printf("Starting with instance name '%s'\n", instanceName)

	// Create necessary dev channels on keybase if they're missing
	wg.Add(1)
	go createMissingChannels(instanceName)

	// Start Rest API
	go apiService()

	// Monitor Keybase Chat API for messages
	kbChannels := []keybase.Channel{
		keybase.Channel{
			Name:      k.Username,
			TopicType: keybase.DEV,
			TopicName: fmt.Sprintf("__%s_input", instanceName),
		},
		keybase.Channel{
			Name:      k.Username,
			TopicType: keybase.DEV,
			TopicName: fmt.Sprintf("__%s_queue", instanceName),
		},
	}
	kbOpts := keybase.RunOptions{
		Dev:            true,
		FilterChannels: kbChannels,
	}
	wg.Wait()
	log.Println("Waiting for new messages...")
	k.Run(handler, kbOpts)
}

func handler(m keybase.ChatAPI) {
	if m.Msg.Content.Type != "text" {
		return
	}
	msg, err := getMessage(m.Msg.Content.Text.Body)
	if err != nil {
		return
	}
	log.Printf("Incoming message ID: %s", msg.ID)

	switch msg.Type {

	case "message":
		// message with type 'message' received
		if event, ok := eventCommands[msg.Event]; ok {
			log.Printf("Message %s triggers eventCommand %s", msg.ID, event.Name)
			event.Exec(msg)
		} else {
			fmt.Println(m.Msg.Content.Text.Body)
		}

		if *msg.Ack {
			// ack message
			log.Printf("Acking message %s", msg.ID)
			var jsonBytes []byte
			jsonBytes, _ = json.Marshal(message{ID: msg.ID, Type: "ack"})
			sendMessage(fmt.Sprintf("__%s_input", msg.Sender), string(jsonBytes))
		}

	case "ack":
		// message with type 'ack' received
		log.Printf("Ack for message %s received", msg.ID)
		repoDestroyMessage(msg.ID)
	}
}
