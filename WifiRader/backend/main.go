package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

type LocationData struct {
    Pos struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    } `json:"pos"`
    DesiredAmount int `json:"desiredAmount"`
}

var (
    data LocationData
    // locMutex        sync.Mutex
    apiKey string
    location string
    radius string
    keyword string
)

type PlaceInfo struct {
    Name      string  `json:"name"`
    URL       string  `json:"url"`
    Latitude  float64 `json:"Latitude"`
    Longitude float64 `json:"Longitude"`
}

func init() {
	// グローバル変数の初期化
	apiKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	radius = "1800" //毎回1500m検索はえぐい
	keyword = "Wi-Fi study"
}

func submitLocationHandler(w http.ResponseWriter, r *http.Request) {
    // コンテンツタイプを最初に設定
    w.Header().Set("Content-Type", "application/json")

    if r.Method != "POST" {
        http.Error(w, `{"error":"Only POST method is allowed"}`, http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading body: %v", err)
        http.Error(w, `{"error":"Can't read body"}`, http.StatusBadRequest)
        return
    }

    if err := json.Unmarshal(body, &data); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusBadRequest)
        return
    }

    log.Printf("Received location: %v, %v and desired amount: %d", data.Pos.Latitude, data.Pos.Longitude, data.DesiredAmount)

    location := fmt.Sprintf("%f,%f", data.Pos.Latitude, data.Pos.Longitude)
    places, err := searchPlaces(apiKey, location, radius, keyword)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to search places: %v", err)
        log.Println(errMsg)
        http.Error(w, fmt.Sprintf(`{"error":"%v"}`, errMsg), http.StatusInternalServerError)
        return
    }

    var response []PlaceInfo
	for _, place := range places.Results {
		url, err := fetchPlaceDetails(apiKey, place.PlaceID)
		if err != nil {
			log.Printf("Failed to fetch details for place %s: %v", place.Name, err)
			continue
		}
		response = append(response, PlaceInfo{Name: place.Name, URL: url, Latitude: place.Geometry.Location.Lat, Longitude: place.Geometry.Location.Lng})
        fmt.Printf("Place: %s, URL: %s\n", place.Name, url)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, `{"error":"Error encoding JSON"}`, http.StatusInternalServerError)
    }
}



func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/submit-location", submitLocationHandler)

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"}, // Node.jsサーバーのポートを指定
        AllowedMethods: []string{"POST"},
        AllowedHeaders: []string{"Content-Type"},
        AllowCredentials: true,
    })

    handler := c.Handler(mux)
    fmt.Println("Server is running on http://localhost:8080")
    http.ListenAndServe(":8080", handler)
}
