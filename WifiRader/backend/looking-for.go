package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type PlaceSearchResponse struct {
	Results []struct {
		Name     string `json:"name"`
		PlaceID  string `json:"place_id"`
		Vicinity string `json:"vicinity"`
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

func searchPlaces(apiKey, location, radius, keyword string) (*PlaceSearchResponse, error) {
	// URLエンコードを使用してキーワードをエンコード
	encodedKeyword := url.QueryEscape(keyword)

	requestURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s&radius=%s&keyword=%s&type=cafe&key=%s", location, radius, encodedKeyword, apiKey)

	resp, err := http.Get(requestURL)
	if err != nil {
		log.Printf("Failed to send request to Google Places API: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// レスポンスボディのデコード
	var result PlaceSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func fetchPlaceDetails(apiKey, placeID string) (string, error) {
	detailURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?placeid=%s&key=%s", placeID, apiKey)
	resp, err := http.Get(detailURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var details struct {
		Result struct {
			Website string `json:"website"`
		} `json:"result"`
		Status string `json:"status"`
	}
	if err := json.Unmarshal(body, &details); err != nil {
		return "", err
	}

	if details.Status != "OK" {
		return "", fmt.Errorf("failed to get place details: %s", details.Status)
	}

	return details.Result.Website, nil
}
