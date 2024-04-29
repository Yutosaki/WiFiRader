package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	// "strings"
	// "time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	// "github.com/gocolly/colly/v2"
)

const (
	targetDomain = "zenn.dev"
	targetURL    = "https://" + targetDomain
)

var (
	endpoint      = os.Getenv("END_POINT")
	subscription  = os.Getenv("ACCOUNT_KEY")
	imageFilePath = ""
	outputFile    = ""
)

// func main() {
// 	c := colly.NewCollector(
// 		// Zenn 以外のアクセスを許可しない
// 		//colly.AllowedDomains(targetDomain),
// 		// ./cache でレスポンスをキャッシュする
// 		colly.CacheDir("./cache"),
// 		// アクセスするページの再帰の深さを設定
// 		//colly.MaxDepth(2),
// 		// ユーザーエージェントを設定
// 		colly.UserAgent("Sample-Scraper"),
// 	)

// 	// リクエスト間で1~2秒の時間を空ける
// 	c.Limit(&colly.LimitRule{
// 		DomainGlob:  targetDomain,
// 		Delay:       time.Second,
// 		RandomDelay: time.Second,
// 	})

// 	// エラー発生時に実行される関数
// 	c.OnError(func(r *colly.Response, err error) {
// 		log.Fatalln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
// 	})
// 	c.OnHTML("img", func(e *colly.HTMLElement) {
// 		imageURL := e.Attr("src")
// 		makeImageFile(imageURL)
// 		fmt.Println("Image URL:", imageURL)
// 		if !strings.HasSuffix(imageURL, "https://movecafe.com/wp/wp-content/uploads/2022/05/sns2.png") {
// 			outputFile = "./output/" + imageURL[len(imageURL)-7:] + ".txt"
// 			ocr()
// 		}
// 	})

// 	url := fmt.Sprintf("%s/topics/go?order=latest", targetURL)
// 	url = "https://cotocafe.jp/menu/"
// 	c.Visit(url)
// }

func ocr() {
	// Create a new Computer Vision client
	client := computervision.New(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscription)

	// Open the image file
	file, err := os.Open(imageFilePath)
	if err != nil {
		log.Fatalf("Error opening image file: %v", err)
	}
	defer file.Close()

	// Perform OCR on the image
	ctx := context.Background()
	result, err := client.RecognizePrintedTextInStream(ctx, true, file, "ja")
	if err != nil {
		log.Fatalf("Error recognizing text: %v", err)
	}

	// Write the recognized text to the output file
	if err := writeTextToFile(result); err != nil {
		log.Fatalf("Error writing text to file: %v", err)
	}

	fmt.Println("Texts are written to output.txt")
}

func makeImageFile(url string) {
	imageFilePath = "./png/" + url[len(url)-7:] + ".png"
	_, err := os.Stat(imageFilePath)
	if !os.IsNotExist(err) {
		err := os.Remove(imageFilePath)
		fmt.Println("IsExist")
		if err != nil {
			log.Fatal(err)
		}
	}
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(imageFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(body)
}

func writeTextToFile(result computervision.OcrResult) error {
	file, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Println("\n************************************")
	for _, region := range *result.Regions {
		for _, line := range *region.Lines {
			for _, word := range *line.Words {
				if word.Text != nil {
					_, err := file.WriteString(*word.Text)
					fmt.Print(*word.Text)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
