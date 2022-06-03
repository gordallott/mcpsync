package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Version struct {
	Version string `json:"version"`
}

func GetVersion(ctx context.Context, targetIP string) (version Version, err error) {
	u, err := url.Parse("http://" + targetIP)
	if err != nil {
		return
	}

	u.Path = "/api/version"

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&version)
	if err != nil {
		return
	}

	return
}
