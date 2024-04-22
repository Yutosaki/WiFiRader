package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "net/url" // この行を追加
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

    // APIリクエストURLの構築
    requestURL := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s&radius=%s&keyword=%s&type=cafe&key=%s", location, radius, encodedKeyword, apiKey)

    // HTTP GETリクエストの送信
    resp, err := http.Get(requestURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    // レスポンスボディの読み込み
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
            URL string `json:"url"`
        } `json:"result"`
        Status string `json:"status"`
    }
    if err := json.Unmarshal(body, &details); err != nil {
        return "", err
    }

    if details.Status != "OK" {
        return "", fmt.Errorf("failed to get place details: %s", details.Status)
    }

    return details.Result.URL, nil
}


func main() {
    apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    location := "35.6895,139.6917" 
    radius := "1500" 
    keyword := "Wi-Fi cafe"

    places, err := searchPlaces(apiKey, location, radius, keyword)
    if err != nil {
        log.Fatalf("Failed to search places: %v", err)
    }

    for _, place := range places.Results {
        url, err := fetchPlaceDetails(apiKey, place.PlaceID)
        if err != nil {
            log.Printf("Failed to fetch details for place %s: %v", place.Name, err)
            continue
        }
        fmt.Printf("Place: %s, URL: %s\n", place.Name, url)
    }
}
