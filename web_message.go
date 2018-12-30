package main

type Message struct {
	Id      string `json:"id"`
	Ack     *bool  `json:"ack",omitempty`
	Target  string `json:"target"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Event   string `json:"event"`
}

type Messages []Message
