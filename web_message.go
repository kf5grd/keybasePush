package main

// A Message struct will be marshaled into JSON data to be sent to the target
// node
type Message struct {
	Id      string `json:"id"`
	Ack     *bool  `json:"ack,omitempty"`
	Target  string `json:"target"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Event   string `json:"event"`
}

// Messages will be a slice of Message structs. This will be used as a queue
// of messages
type Messages []Message
