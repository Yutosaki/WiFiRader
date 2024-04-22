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

func init() {
	// グローバル変数の初期化
	apiKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	radius = "1500"
	keyword = "Wi-Fi cafe"
}

func submitLocationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading body: %v", err)
        http.Error(w, "can't read body", http.StatusBadRequest)
        return
    }

    // log.Printf("Body received: %s", body)

    
    if err := json.Unmarshal(body, &data); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received location: %v, %v and desired amount: %d\n", data.Pos.Latitude, data.Pos.Longitude, data.DesiredAmount)
    fmt.Fprintf(w, "Data received successfully")

    location := fmt.Sprintf("%f,%f", data.Pos.Latitude, data.Pos.Longitude)
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

    // apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
    // location := fmt.Sprintf("%f,%f", currentLocation.Pos.Latitude, currentLocation.Pos.Longitude)
    // radius := "1500" 
    // keyword := "Wi-Fi cafe"

    // places, err := searchPlaces(apiKey, location, radius, keyword)
    // if err != nil {
    //     log.Fatalf("Failed to search places: %v", err)
    // }

    // for _, place := range places.Results {
    //     url, err := fetchPlaceDetails(apiKey, place.PlaceID)
    //     if err != nil {
    //         log.Printf("Failed to fetch details for place %s: %v", place.Name, err)
    //         continue
    //     }
    //     fmt.Printf("Place: %s, URL: %s\n", place.Name, url)
    // }
}
