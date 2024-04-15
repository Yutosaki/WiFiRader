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
        Vicinity string `json:"vicinity"`
    } `json:"results"`
    Status string `json:"status"`
}

func searchPlaces(apiKey, location, radius, keyword string) (*PlaceSearchResponse, error) {
    // キーワードをURLエンコード
    encodedKeyword := url.QueryEscape(keyword)

    url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s&radius=%s&keyword=%s&type=cafe&key=%s", location, radius, encodedKeyword, apiKey)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result PlaceSearchResponse
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    return &result, nil
}

func main() {
    apiKey := "AIzaSyAOOq8iYHO0xynkfOnUfHOm00VjO72Wufk"
    location := "35.6895,139.6917" // 東京の座標、実際にはユーザーの現在地を使用
    radius := "1500" // 単位はメートル
    keyword := "Wi-Fi cafe study"

    places, err := searchPlaces(apiKey, location, radius, keyword)
    if err != nil {
        log.Fatalf("Failed to search places: %v", err)
    }
    fmt.Printf("Places found: %+v\n", places)
}
