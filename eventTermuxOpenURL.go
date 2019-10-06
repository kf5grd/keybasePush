// +build linux,arm64

package main

import (
	"log"
	"os/exec"
)

func init() {
	var command = eventCommand{
		Name:   "TermuxOpenURL",
		Events: []string{"openurl"},
		Exec:   eventTermuxOpenURL,
	}

	registerEventCommand(command)
}

func eventTermuxOpenURL(m message) {
	_, err := exec.Command("termux-open-url", m.Content).Output()
	if err != nil {
		log.Printf("Unable to open url from message %s: %s", m.ID, err)
	}
}
