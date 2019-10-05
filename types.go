package main

// A message struct will be marshaled into JSON data to be sent to the target
// node
type message struct {
	ID      string `json:"id"`
	Type    string `json:"type"` // 'message' or 'ack'
	Ack     *bool  `json:"ack,omitempty"`
	Sender  string `json:"sender,omitempty"`
	Target  string `json:"target,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Event   string `json:"event,omitempty"`
}

// eventCommand holds a command that can be triggered by messages with certain events
type eventCommand struct {
	Name   string        // name to use when referencing this command in logs
	Events []string      // events that will trigger this command
	Exec   func(message) // function to execute when this command is triggered
}
