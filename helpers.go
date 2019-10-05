package main

import (
	"encoding/json"
	"fmt"
	"log"

	"samhofi.us/x/keybase"
)

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

// sendMessage will send a string to a specified channel
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
