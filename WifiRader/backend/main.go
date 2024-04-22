package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"github.com/rs/cors"
)

type LocationData struct {
    Pos struct {
        Latitude  float64 `json:"latitude"`
        Longitude float64 `json:"longitude"`
    } `json:"pos"`
    DesiredAmount int `json:"desiredAmount"`
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

    log.Printf("Body received: %s", body)

    var data LocationData
    if err := json.Unmarshal(body, &data); err != nil {
        log.Printf("Error decoding JSON: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received location: %v, %v and desired amount: %d\n", data.Pos.Latitude, data.Pos.Longitude, data.DesiredAmount)
    fmt.Fprintf(w, "Data received successfully")
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
