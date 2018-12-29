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
	Channel   chatOutChannel `json:"channel,omitempty"`
	Message   chatOutMessage `json:"message,omitempty"`
	TopicType string         `json:"topic_type,omitempty"`
}

type chatOutChannel struct {
	Name        string `json:"name"`
	MembersType string `json:"members_type,omitempty"`
	TopicName   string `json:"topic_name,omitempty"`
	TopicType   string `json:"topic_type,omitempty"`
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
	Message       string               `json:"message"`
	Conversations []chatInConversation `json:"conversations"`
	RateLimits    []chatInLimits       `json:"ratelimits"`
}

type chatInLimits struct {
	Tank     string `json:"tank"`
	Capacity int    `json:"capacity"`
	Reset    int    `json:"reset"`
	Gas      int    `json:"gas"`
}

type chatInConversation struct {
	Channel conversationChannel `json:"channel"`
}

type conversationChannel struct {
	Name      string `json:"name"`
	TopicType string `json:"topic_type"`
	TopicName string `json:"topic_name"`
}

type chatInError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendChatAPI(jsonData string) (chatAPIIn, error) {
	cmd := exec.Command("keybase", "chat", "api", "-m", jsonData)

	cmdOut, err := cmd.Output()
	if err != nil {
		return chatAPIIn{}, err
	}

	var retVal chatAPIIn
	json.Unmarshal(cmdOut, &retVal)

	if retVal.Error.Message != "" {
		return chatAPIIn{}, errors.New(retVal.Error.Message)
	}

	return retVal, nil
}

func SendDevMessage(user, channelName, message string) error {
	var msgJSON = chatAPIOut{
		Method: "send",
		Params: chatOutParams{
			Options: chatOutOptions{
				Channel: chatOutChannel{
					Name:      user,
					TopicName: channelName,
					TopicType: "DEV",
				},
				Message: chatOutMessage{
					Body: message,
				},
			},
		},
	}
	jsonBytes, _ := json.Marshal(msgJSON)
	_, err := SendChatAPI(string(jsonBytes))
	if err != nil {
		return err
	}
	return nil
}

func GetDevChannels() []string {
	var msgJSON = chatAPIOut{
		Method: "list",
		Params: chatOutParams{
			Options: chatOutOptions{
				TopicType: "DEV",
			},
		},
	}

	jsonBytes, _ := json.Marshal(msgJSON)
	allChannels, err := SendChatAPI(string(jsonBytes))
	if err != nil {
		panic(err)
	}

	devChannels := []string{}
	user := KeybaseUsername()
	for _, channel := range allChannels.Result.Conversations {
		if (channel.Channel.Name == user) && (channel.Channel.TopicType == "dev") {
			devChannels = append(devChannels, channel.Channel.TopicName)
		}
	}
	return devChannels
}

func CreateDevChannel(user, channelName string) error {
	cmd := exec.Command("keybase", "chat", "send", "--topic-type", "dev", "--channel", channelName, user, "{\"create_channel\": true}")

	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}
