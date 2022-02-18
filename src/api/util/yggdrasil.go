package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type AuthenticationResponse struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Properties []struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		Signature string `json:"signature"`
	} `json:"properties"`
}

func Authenticate(username, hash string) (*AuthenticationResponse, error) {
	resp, err := http.Get("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=" + username + "&serverId=" + hash)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var response AuthenticationResponse

	err = json.Unmarshal(body, &response)

	return &response, err
}
