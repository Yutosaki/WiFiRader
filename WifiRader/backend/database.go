package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

type PlaceInfo struct {
	Name         string
	URL          string
	Latitude     float64
	Longitude    float64
	LastChecked  time.Time
	MinimumPrice int
}

var db *sql.DB

func initDB() error {
	// .envファイルの読み込み
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// 環境変数からデータベース接続情報を取得
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// データソース名（DSN）を作成
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port)
	fmt.Println("Connecting to database with DSN:", dsn)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// データベースを作成
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbname)
	if err != nil {
		return fmt.Errorf("error creating database: %v", err)
	}

	// データベースを選択
	_, err = db.Exec("USE " + dbname)
	if err != nil {
		return fmt.Errorf("error selecting database: %v", err)
	}

	// テーブル作成スクリプトを実行
	err = executeSQLScript("table.sql")
	if err != nil {
		return fmt.Errorf("error executing table creation script: %v", err)
	}

	return nil
}

func executeSQLScript(filename string) error {
	script, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading SQL script file: %v", err)
	}

	// SQLコマンドを実行する
	_, err = db.Exec(string(script))
	if err != nil {
		return fmt.Errorf("error executing SQL script: %v", err)
	}

	return nil
}

func upsertPlace(place PlaceInfo) error {
	today := time.Now().Format("2006-01-02")

	// Placesテーブルに存在するか確認
	var placeID int
	err := db.QueryRow("SELECT PlaceID FROM Places WHERE URL = ?", place.URL).Scan(&placeID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Placeが存在しない場合は新規挿入
			res, err := db.Exec("INSERT INTO Places (Name, URL, Latitude, Longitude) VALUES (?, ?, ?, ?)", place.Name, place.URL, place.Latitude, place.Longitude)
			if err != nil {
				return fmt.Errorf("error inserting place: %v", err)
			}

			placeID64, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("error getting last insert id: %v", err)
			}
			placeID = int(placeID64)
		} else {
			return fmt.Errorf("error querying place: %v", err)
		}
	}

	// Pricesテーブルに挿入または更新
	var lastChecked []byte
	err = db.QueryRow("SELECT LastChecked FROM Prices WHERE PlaceID = ? ORDER BY LastChecked DESC LIMIT 1", placeID).Scan(&lastChecked)
	if err != nil {
		if err == sql.ErrNoRows || parseDate(lastChecked).AddDate(0, 1, 0).Before(time.Now()) {
			// 新規挿入または1ヶ月以上前のデータの場合は更新
			_, err = db.Exec("INSERT INTO Prices (PlaceID, LastChecked, MinimumPrice) VALUES (?, ?, ?)", placeID, today, place.MinimumPrice)
			if err != nil {
				return fmt.Errorf("error inserting/updating price: %v", err)
			}
		} else if err != nil {
			return fmt.Errorf("error querying price: %v", err)
		}
	}
	return nil
}

// parseDateはバイトスライスからtime.Timeに変換します
func parseDate(dateBytes []byte) time.Time {
	dateStr := string(dateBytes)
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Printf("Error parsing date: %v", err)
		return time.Time{}
	}
	return date
}

// getPlaceInfoはURLに基づいてPlace情報を取得
func getPlaceInfo(url string) (*PlaceInfo, error) {
	var place PlaceInfo
	var lastChecked []byte
	err := db.QueryRow(`
		SELECT p.Name, p.URL, p.Latitude, p.Longitude, pr.MinimumPrice, pr.LastChecked 
		FROM Places p
		JOIN Prices pr ON p.PlaceID = pr.PlaceID
		WHERE p.URL = ?
		ORDER BY pr.LastChecked DESC
		LIMIT 1`, url).Scan(&place.Name, &place.URL, &place.Latitude, &place.Longitude, &place.MinimumPrice, &lastChecked)
	if err != nil {
		return nil, fmt.Errorf("error querying place info: %v", err)
	}
	place.LastChecked = parseDate(lastChecked)
	return &place, nil
}

func main() {
	// データベース初期化
	err := initDB()
	if err != nil {
		log.Fatalf("DB initialization failed: %v", err)
	}
	defer db.Close()

	// サンプルデータの挿入
	place := PlaceInfo{
		Name:         "Example Cafe",
		URL:          "https://examplecafe.com",
		Latitude:     35.6895,
		Longitude:    139.6917,
		MinimumPrice: 500,
	}

	err = upsertPlace(place)
	if err != nil {
		log.Fatalf("Failed to upsert place: %v", err)
	}

	// データの取得
	retrievedPlace, err := getPlaceInfo("https://examplecafe.com")
	if err != nil {
		log.Fatalf("Failed to get place info: %v", err)
	}

	fmt.Printf("Retrieved Place: %+v\n", retrievedPlace)
}
