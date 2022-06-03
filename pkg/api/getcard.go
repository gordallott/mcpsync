package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func GetCardFilename(gameID string, channel uint) string {
	return fmt.Sprintf("%s-%d.mcd", gameID, channel)
}

func GetCard(ctx context.Context, targetIP, gameID string, channel uint) (io.ReadCloser, error) {
	if channel < 1 || channel > 8 {
		return nil, fmt.Errorf("channel must be between 1 and 8")
	}

	if gameID == "" {
		return nil, fmt.Errorf("gameID must be set")
	}

	u, err := url.Parse("http://" + targetIP)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join("MemoryCards", gameID, GetCardFilename(gameID, channel))

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if statusCode := response.StatusCode; statusCode != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", statusCode)
	}

	return response.Body, nil
}
