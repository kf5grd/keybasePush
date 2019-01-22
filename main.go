package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

const (
	// IP address to serve rest API on
	serverHost string = "127.0.0.1"
)

// Initialize variables for flags
var instanceName string
var serverPort int

func init() {
	// Parse flags
	flag.StringVar(&instanceName, "name", "", "Set the name of this instance")
	flag.IntVar(&serverPort, "port", 8617, "Set the port for the API")
	flag.Parse()
}

func apiService() {
	// Create a new router and start the rest API service
	router := NewRouter()
	log.Printf("Starting with instance name '%s'\n", instanceName)
	serverAddress := fmt.Sprintf("%s:%d", serverHost, serverPort)
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

func main() {
	// If -name flag isn't set, use device name from keybase
	if instanceName == "" {
		instanceName = KeybaseDeviceName()
	}
	// Instance name forced to lowercase with no leading or trailing spaces
	instanceName = strings.ToLower(instanceName)
	instanceName = strings.TrimSpace(instanceName)
	// Create necessary dev channels on keybase if they're missing
	CreateMissingChannels(instanceName)

	go apiService()

	// Monitor Keybase Chat API for messages
	listener := exec.Command("keybase", "chat", "api-listen", "--dev")
	listenerOutput, _ := listener.StdoutPipe()
	listener.Start()
	scanner := bufio.NewScanner(listenerOutput)
	log.Println("Waiting for new messages...")
	for scanner.Scan() {
		newMessage := ReceiveMessage(scanner.Text())
		msgSender := newMessage.Msg.Sender.Username
		contentType := newMessage.Msg.Content.Type
		topicType := newMessage.Msg.Channel.TopicType
		channelName := newMessage.Msg.Channel.TopicName

		if (msgSender == KeybaseUsername()) && (topicType == "dev") && (channelName == fmt.Sprintf("__%s_input", instanceName)) && (contentType == "text") {
			emptyMessage := Message{}
			if msg := GetMessage(newMessage.Msg.Content.Text.Body); msg != emptyMessage {
				switch msg.Type {

				case "message":
					// message with type 'message' received
					fmt.Println(newMessage.Msg.Content.Text.Body)
					if *msg.Ack {
						// ack message
						var jsonBytes []byte
						jsonBytes, _ = json.Marshal(Message{Id: msg.Id, Type: "ack"})
						SendDevMessage(KeybaseUsername(), fmt.Sprintf("__%s_input", msg.Sender), string(jsonBytes))
					}

				case "ack":
					// message with type 'ack' received
				}
			}
		}
	}
}

// Get content from received message
func GetMessage(jsonString string) Message {
	var jsonData Message
	json.Unmarshal([]byte(jsonString), &jsonData)
	return jsonData
}

// CreateMissingChannels will check if the queue and input 'dev' channels have
// already been created on keybase, and if not, it will create them
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
