package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Card struct {
	FullPath string `json:"fullpath"`
	GameID   string `json:"gameId"`
	Name     string `json:"name"`
}

type CardsPaginated struct {
	Total   int    `json:"total"`
	Results []Card `json:"results"`
}

func getCardsPage(ctx context.Context, targetIP string, start, limit uint) (cards CardsPaginated, err error) {
	u, err := url.Parse("http://" + targetIP)
	if err != nil {
		return
	}

	u.Path = "/api/query"

	type CardsPaginatedReq struct {
		Start uint `json:"start"`
		Limit uint `json:"limit"`
	}

	reqBytes, err := json.Marshal(CardsPaginatedReq{
		Start: start,
		Limit: limit,
	})

	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), bytes.NewReader(reqBytes))
	if err != nil {
		return
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&cards)
	if err != nil {
		return
	}

	return
}

func GetCards(ctx context.Context, targetIP string) (cards []Card, err error) {
	limit := uint(10)
	for start := uint(0); ; start += limit {
		cardsPage, err := getCardsPage(ctx, targetIP, start, limit)
		if err != nil {
			return nil, err
		}

		cards = append(cards, cardsPage.Results...)

		if len(cardsPage.Results) < int(limit) {
			break
		}
	}

	return
}
