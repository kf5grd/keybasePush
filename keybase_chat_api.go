package main

import (
	"encoding/json"
	"errors"
	"os/exec"
)

// JSON Out to API
type chatAPIOut struct {
	Method string        `json:"method"`
	Params chatOutParams `json:"params,omitempty"`
}

type chatOutParams struct {
	Options chatOutOptions `json:"options,omitempty"`
}

type chatOutOptions struct {
	Channel chatOutChannel `json:"channel,omitempty"`
	Message chatOutMessage `json:"message,omitempty"`
}

type chatOutChannel struct {
	Name        string `json:"name"`
	MembersType string `json:"members_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
}

type chatOutMessage struct {
	Body string `json:"body"`
}

// JSON Received back from API
type chatAPIIn struct {
	Result chatInResult `json:"result,omitempty"`
	Error  chatInError  `json:"error,omitempty"`
}

type chatInResult struct {
	Message    string         `json:"message"`
	RateLimits []chatInLimits `json:"ratelimits"`
}

type chatInLimits struct {
	Tank     string `json:"tank"`
	Capacity int    `json:"capacity"`
	Reset    int    `json:"reset"`
	Gas      int    `json:"gas"`
}

type chatInError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendChatAPI(jsonData string) (chatAPIIn, error) {
	cmd := exec.Command("keybase", "chat", "api", "-m", jsonData)

	cmdOut, err := cmd.Output()
	if err != nil {
		panic(err)
	}

	var retVal chatAPIIn
	json.Unmarshal(cmdOut, &retVal)

	if retVal.Error.Message != "" {
		return chatAPIIn{}, errors.New(retVal.Error.Message)
	}

	return retVal, nil
}