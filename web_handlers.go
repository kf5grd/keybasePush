package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Path: /
// Method: GET
// Purpose: Currently displays instance name
// ToDo: Add some instructions for using the system. Maybe in the distant
//       future this page can display a form for making some config changes.
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to %s!", instanceName)
}

// Path: /messages
// Method: GET
// Purpose: Displays full message queue. This shouldn't be needed as the queue
//          is tracked within keybase and external processes shouldn't need to
//          access the queue.
func messageIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		panic(err)
	}
}

// Path: /messages/<messageId>
// Method: GET
// Purpose: Get an individual message from the queue. This shouldn't be needed as the queue
//          is tracked within keybase and external processes shouldn't need to
//          access the queue.
func messageShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	msgID := vars["messageId"]
	var msg message
	msg = repoFindMessage(msgID)

	if msg.ID == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Message not found"}); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		panic(err)
	}
}

// Path: /messages
// Method: POST
// Purpose: Sends a JSON encoded message object to target instance and adds
//          message to queue.
func MessageCreate(w http.ResponseWriter, r *http.Request) {
	var msg message
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessible entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	if msg.Content == "" || msg.Target == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessible entity

		resp := map[string]string{"Error": "'target' and 'content' fields are required"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	// Ack defaults to true
	if msg.Ack == nil {
		t := new(bool)
		*t = true
		msg.Ack = t
	}

	// Set msg type
	msg.Type = "message"

	// Force target to lowercase
	msg.Target = strings.ToLower(msg.Target)

	// Attach instance name to msg
	msg.Sender = strings.ToLower(instanceName)

	m := repoCreateMessage(msg)

	// Send msg to input channel for target
	channel := fmt.Sprintf("__%s_input", m.Target)
	jsonBytes, _ := json.Marshal(m)
	if err := sendMessage(channel, string(jsonBytes)); err != nil {
		// delete msg
		repoDestroyMessage(m.ID)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Message not delivered"}); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(m); err != nil {
		panic(err)
	}
}
