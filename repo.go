package main

import (
	"crypto/sha1"
	"fmt"
	"time"
)

var messageId string
var messages Messages

func SendQueue(m Messages) error {
	// Send updated queue
	channel = fmt.Sprintf("__%s_queue", instanceName)
	jsonBytes, _ = json.Marshal(m)
	if err := SendDevMessage(user, channel, string(jsonBytes)); err != nil {
		return err
	}
	return nil
}

func RepoFindMessage(id string) Message {
	for _, m := range messages {
		if m.Id == id {
			return m
		}
	}
	// return empty message if not found
	return Message{}
}

func RepoCreateMessage(m Message) Message {
	data := []byte(fmt.Sprintf("%s%s", time.Now(), m.Content))
	currentId := fmt.Sprintf("%x", sha1.Sum(data))[:8]
	m.Id = currentId

	newMessages = append(messages, m)
	if err := SendQueue(); err != nil {
		emptyMessage := Message{}
		return emptyMessage
	}

	message = newMessages
	return m
}

func RepoDestroyMessage(id string) error {
	for i, m := range messages {
		if m.Id == id {
			newMessages = append(messages[:i], messages[i+1:]...)
			if err := SendQueue(); err != nil {
				return err
			}
			messages = newMessages
			return nil
		}
	}
	return fmt.Errorf("Could not find Message with id %s to delete", id)
}
