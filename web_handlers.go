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

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to %s!", instanceName)
}

func MessageIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		panic(err)
	}
}

func MessageShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	var message Message
	message = RepoFindMessage(messageId)

	if message.Id == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Message not found"}); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(message); err != nil {
		panic(err)
	}
}

func MessageCreate(w http.ResponseWriter, r *http.Request) {
	var message Message
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &message); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessible entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	if message.Content == "" || message.Target == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessible entity

		resp := map[string]string{"Error": "'target' and 'content' fields are required"}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	// Ack defaults to true
	if message.Ack == nil {
		t := new(bool)
		*t = true
		message.Ack = t
	}

	// Force target to lowercase
	message.Target = strings.ToLower(message.Target)

	m := RepoCreateMessage(message)

	// Send message to input channel for each target
	user := KeybaseUsername()
	channel := fmt.Sprintf("__%s_input", m.Target)
	jsonBytes, _ := json.Marshal(m)
	if err := SendDevMessage(user, channel, string(jsonBytes)); err != nil {
		// delete message
		RepoDestroyMessage(m.Id)

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

	// Send updated queue
	channel = fmt.Sprintf("__%s_queue", instanceName)
	jsonBytes, _ = json.Marshal(messages)
	SendDevMessage(user, channel, string(jsonBytes))
}
