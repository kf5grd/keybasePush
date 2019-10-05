package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func init() {
	// Parse flags
	flag.StringVar(&instanceName, "name", "", "Set the name of this instance")
	flag.IntVar(&serverPort, "port", 8617, "Set the port for the API")
	flag.Parse()
}

func apiService() {
	// Create a new router and start the rest API service
	router := newRouter()
	log.Printf("Starting with instance name '%s'\n", instanceName)
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

	// Create necessary dev channels on keybase if they're missing
	createMissingChannels(instanceName)

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

	chat := k.NewChat(m.Msg.Channel)
	switch msg.Type {

	case "message":
		// message with type 'message' received
		fmt.Println(m.Msg.Content.Text.Body)
		if *msg.Ack {
			// ack message
			var jsonBytes []byte
			jsonBytes, _ = json.Marshal(message{ID: msg.ID, Type: "ack"})
			chat.Send(string(jsonBytes))
		}

	case "ack":
		// message with type 'ack' received
		repoDestroyMessage(msg.ID)
	}
}

// getMessage gets content from received message
func getMessage(jsonString string) (message, error) {
	var result message
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// createMissingChannels will check if the queue and input 'dev' channels have
// already been created on keybase, and if not, it will create them
func createMissingChannels(instanceName string) {
	// need queue and input channels
	neededChannels := []string{
		fmt.Sprintf("__%s_queue", instanceName),
		fmt.Sprintf("__%s_input", instanceName),
	}

	// get existing dev channels
	existingChannels := make(map[string]bool)
	devChannels, err := k.ChatList(keybase.Channel{
		Name:      k.Username,
		TopicType: keybase.DEV,
	})
	if err != nil {
		log.Printf("Unable to get dev channels: %v", err)
	} else {
		for _, c := range devChannels.Result.Conversations {
			if c.Channel.Name == k.Username {
				existingChannels[c.Channel.TopicName] = true
			}
		}
	}

	// create any dev channel that's needed and doesn't exist yet
	for _, devChan := range neededChannels {
		if _, ok := existingChannels[devChan]; !ok {
			log.Printf("Creating missing dev channel '%s'\n", devChan)
			chat := k.NewChat(keybase.Channel{
				Name:        k.Username,
				MembersType: keybase.USER,
				TopicName:   devChan,
				TopicType:   keybase.DEV,
			})
			_, err := chat.Send(`{"create_channel": true}`)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func sendMessage(channel string, msg string) error {
	ch := keybase.Channel{
		Name:        k.Username,
		MembersType: keybase.USER,
		TopicName:   channel,
		TopicType:   keybase.DEV,
	}
	chat := k.NewChat(ch)
	_, err := chat.Send(msg)
	return err
}
