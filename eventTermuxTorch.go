// +build linux,arm64

package main

import (
	"log"
	"os/exec"
	"strings"
)

func init() {
	var command = eventCommand{
		Name:   "TermuxTorch",
		Events: []string{"torch"},
		Exec:   eventTermuxTorch,
	}

	registerEventCommand(command)
}

func eventTermuxTorch(m message) {
	c := strings.ToLower(m.Content)
	c = strings.Trim(c, " ")
	if c == "on" || c == "off" {
		_, err := exec.Command("termux-torch", c).Output()
		if err != nil {
			log.Printf("Unable to send notification for message %s: %s", m.ID, err)
		}
	}
}
