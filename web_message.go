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
