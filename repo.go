package main

import (
	"crypto/sha1"
	"fmt"
	"time"
)

var messageId string
var messages Messages

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

	messages = append(messages, m)
	return m
}

func RepoDestroyMessage(id string) error {
	for i, m := range messages {
		if m.Id == id {
			messages = append(messages[:i], messages[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Message with id %s to delete", id)
}
