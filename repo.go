package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"time"
)

var messageID string
var messages []message

func sendQueue(m []message) error {
	// Send updated queue
	channel := fmt.Sprintf("__%s_queue", instanceName)
	jsonBytes, _ := json.Marshal(m)
	if err := sendMessage(channel, string(jsonBytes)); err != nil {
		return err
	}
	return nil
}

func repoFindMessage(id string) message {
	for _, m := range messages {
		if m.ID == id {
			return m
		}
	}
	// return empty message if not found
	return message{}
}

func repoCreateMessage(m message) message {
	data := []byte(fmt.Sprintf("%s%s", time.Now(), m.Content))
	currentID := fmt.Sprintf("%x", sha1.Sum(data))[:8]
	m.ID = currentID

	newMessages := append(messages, m)
	if err := sendQueue(newMessages); err != nil {
		emptyMessage := message{}
		return emptyMessage
	}

	messages = newMessages
	return m
}

func repoDestroyMessage(id string) error {
	for i, m := range messages {
		if m.ID == id {
			newMessages := append(messages[:i], messages[i+1:]...)
			if err := sendQueue(newMessages); err != nil {
				return err
			}
			messages = newMessages
			return nil
		}
	}
	return fmt.Errorf("Could not find Message with id %s to delete", id)
}
