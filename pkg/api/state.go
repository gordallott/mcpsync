package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type State struct {
	GameName       string `json:"gameName"`
	GameID         string `json:"gameId"`
	CurrentChannel uint   `json:"currentChannel"`
	RSSI           int    `json:"rssi"`
}

func GetState(ctx context.Context, targetIP string) (state State, err error) {
	u, err := url.Parse("http://" + targetIP)
	if err != nil {
		return
	}

	u.Path = "/api/currentState"

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&state)
	if err != nil {
		return
	}

	return

}
