package main

import (
    "encoding/json"
    "fmt"
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

    var data LocationData
    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    log.Printf("Received location: %v, %v and desired amount: %d\n", data.Pos.Latitude, data.Pos.Longitude, data.DesiredAmount)
    // ここでデータを処理する

    fmt.Fprintf(w, "Data received successfully")
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/submit-location", submitLocationHandler)

    // CORS設定: 全てのオリジンからのリクエストを許可する場合（必要に応じてセキュリティ設定を調整）
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowCredentials: true,
        AllowedHeaders: []string{"Content-Type"},
    })

    handler := c.Handler(mux)
    http.ListenAndServe(":8080", handler)
}
