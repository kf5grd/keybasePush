// +build linux,arm64

package main

import (
	"log"
	"os/exec"
)

func init() {
	registerEventCommand(command)
}

var command = eventCommand{
	Name:   "TermuxNotify",
	Events: []string{"notify"},
	Exec:   eventTermuxNotify,
}

func eventTermuxNotify(m message) {
	_, err := exec.Command("termux-notification", "--title", m.Title, "--content", m.Content).Output()
	if err != nil {
		log.Printf("Unable to send notification for message %s: %s", m.ID, err)
	}
}
