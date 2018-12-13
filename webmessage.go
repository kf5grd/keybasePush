package main

type Message struct {
	Id      string   `json:"id"`
	Targets []string `json:"targets"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Event   string   `json:"event"`
}

type Messages []Message
