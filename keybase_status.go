package main

import (
	"encoding/json"
	"os/exec"
)

var KeybaseStatus keybaseStatus

func init() {
	var err error
	KeybaseStatus, err = GetKeybaseStatus()
	if err != nil {
		panic(err)
	}
}

type keybaseStatus struct {
	Username string        `json:"Username"`
	LoggedIn bool          `json:"LoggedIn"`
	Device   keybaseDevice `json:"Device"`
}

type keybaseDevice struct {
	Name string `json:"name"`
}

// Parse the keybase status command
func GetKeybaseStatus() (keybaseStatus, error) {
	cmd := exec.Command("keybase", "status", "-j")

	cmdOut, err := cmd.Output()
	if err != nil {
		return keybaseStatus{}, err
	}

	var retVal keybaseStatus
	json.Unmarshal(cmdOut, &retVal)

	return retVal, nil
}

// Return the local device name as it shows in the keybase status command
func KeybaseDeviceName() string {
	return KeybaseStatus.Device.Name
}

// Return the logged in keybase user's username as it shows in the keybase
// status command
func KeybaseUsername() string {
	return KeybaseStatus.Username
}

// Return true if keybase client is logged in, otherwise return false
func KeybaseLoggedIn() bool {
	return KeybaseStatus.LoggedIn
}
