package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	"github.com/gocolly/colly/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var (
	endpoint                 = os.Getenv("END_POINT")
	subscription             = os.Getenv("ACCOUNT_KEY")
	geminiapikey             = os.Getenv("GEMINI_API_KEY")
	genaiTrue     genai.Text = "True"
	isTrue                   = false
	visited                  = make(map[string]bool)
	imageFilePath            = ""
	maxprice      int
)

func checkmenu(price int, resp []PlaceInfo) (response []PlaceInfo) {
	os.Mkdir("png", 0777)
	maxprice = price
	var i = 0
	for i < len(resp) {
		isTrue = false
		url := resp[i].URL
		if visited[url] || scraping(url) {
			visited[url] = true
			response = append(response, resp[i])
		}
		os.RemoveAll("png")
		i++
	}
	return response
}

func scraping(url string) bool {
	c := colly.NewCollector(
		// Zenn 以外のアクセスを許可しない
		//colly.AllowedDomains(targetDomain),
		// ./cache でレスポンスをキャッシュする
		//colly.CacheDir("./cache"),
		// アクセスするページの再帰の深さを設定
		//colly.MaxDepth(2),
		// ユーザーエージェントを設定
		colly.UserAgent("Sample-Scraper"),
	)

	// 	// リクエスト間で1~2秒の時間を空ける
	// 	c.Limit(&colly.LimitRule{
	// 		DomainGlob:  targetDomain,
	// 		Delay:       time.Second,
	// 		RandomDelay: time.Second,
	// 	})

	// エラー発生時に実行される関数
	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		prompt := []genai.Part{
			genai.Text(strings.Join(strings.Fields(e.Text), " ")),
		}
		if geminiChat(prompt) {
			isTrue = true
		}
	})
	c.OnHTML("img", func(e *colly.HTMLElement) {
		imageURL := e.Attr("src")
		if !visited[imageURL] {
			if len(imageURL) > 10 && imageURL[4:] == "http" {
				log.Printf("imageURL")
				makeImageFile(imageURL)
				time.Sleep(3 * time.Second)
				if ocr() {
					isTrue = true
				}

			}
			visited[imageURL] = true
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		menuURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Contains(menuURL, "menu") && menuURL != url && !visited[url] {
			if scraping(menuURL) {
				isTrue = true
			}
		}
		visited[url] = true
	})

	c.Visit(url)
	return isTrue
}

func ocr() bool {
	// Create a new Computer Vision clien
	client := computervision.New(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscription)

	// Open the image file
	file, err := os.Open(imageFilePath)
	if err != nil {
		//log.Fatalf("Error opening image file: %v", err)
		log.Printf("Error opening image file: %v\n", err)
		return false
	}
	defer file.Close()

	fileinfo, staterr := file.Stat()

	if staterr != nil {
		log.Printf("staterr error: %v\n", staterr)
		return false
	}
	if fileinfo.Size() < 100000 {
		// log.Printf("filesize is so small")
		return false
	}

	// Perform OCR on the image
	ctx := context.Background()
	result, err := client.RecognizePrintedTextInStream(ctx, true, file, "ja")
	if err != nil {
		//log.Fatalf("Error recognizing text: %v", err)
		log.Printf("Error recognizing text: %v", err)
		return false
	}

	// Write the recognized text to the output file
	if err, ok := writeTextToFile(result); err != nil {
		//log.Fatalf("Error writing text to file: %v", err)
		log.Printf("Error writing text to file: %v", err)
		return false
	} else {
		return ok
	}
}

func makeImageFile(url string) {
	imageFilePath = "./png/" + url[len(url)-7:] + ".png"
	_, err := os.Stat(imageFilePath)
	if !os.IsNotExist(err) {
		err := os.Remove(imageFilePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error get url: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error read all: %v", err)
	}

	file, err := os.Create(imageFilePath)
	if err != nil {
		log.Printf("Error create image file: %v", err)
	}
	defer file.Close()

	file.Write(body)
}

func writeTextToFile(result computervision.OcrResult) (error, bool) {
	var prompt []genai.Part
	for _, region := range *result.Regions {
		for _, line := range *region.Lines {
			for _, word := range *line.Words {
				if word.Text != nil {
					prompt = append(prompt, genai.Text(*word.Text))
				}
			}
		}
	}
	return nil, geminiChat(prompt)
}

func geminiChat(prompt []genai.Part) bool {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiapikey))
	if err != nil {
		log.Print(prompt)
		log.Printf("Error generate new client:%v", err)
		return false
	}
	model := client.GenerativeModel("gemini-pro")
	prompt = append(prompt, genai.Text("What is the lowest price. Answer true or false whether it is less than"))
	prompt = append(prompt, genai.Text(strconv.Itoa(maxprice)))
	prompt = append(prompt, genai.Text("Don't say anything other than True or False."))
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Printf("Error generate content:%v", err)
		return false
	}
	return printResponse(resp)
}
func printResponse(resp *genai.GenerateContentResponse) bool {
	for _, candidate := range resp.Candidates {
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				if part == genaiTrue {
					return true
				}
			}
		}
	}
	return false
}
