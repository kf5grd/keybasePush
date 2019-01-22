package main

// A Message struct will be marshaled into JSON data to be sent to the target
// node
type Message struct {
	Id      string `json:"id"`
	Type    string `json:"type"` // 'message' or 'ack'
	Ack     *bool  `json:"ack,omitempty"`
	Sender  string `json:"sender,omitempty"`
	Target  string `json:"target,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	Event   string `json:"event,omitempty"`
}

// Messages will be a slice of Message structs. This will be used as a queue
// of messages
type Messages []Message
