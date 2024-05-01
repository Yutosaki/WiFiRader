package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v3.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	"github.com/gocolly/colly/v2"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	targetDomain = "zenn.dev"
	targetURL    = "https://" + targetDomain
)

var (
	endpoint      = os.Getenv("END_POINT")
	subscription  = os.Getenv("ACCOUNT_KEY")
	geminiapikey  = os.Getenv("GEMINI_API_KEY")
	genaiTrue genai.Text = "True"
	imageFilePath = ""
	outputFile    = ""
)

var maxprice = 49

func main() {
	os.Mkdir("png", 0777)
	os.Mkdir("output", 0777)
	url := "http://nericafe.com/"
	url = "https://www.tullys.co.jp/menu/drink/coffee/coldbrewcoffee24.html"
	url = "https://cotocafe.jp/menu/"
	//url = "https://bowlscafe.com/menu"
	scraping(url)
	if maxprice != 0 {
		fmt.Println(maxprice)
	}
	fmt.Println("if error occured when remove png dir:", os.RemoveAll("png"))
	fmt.Println("if error occured when remove output dir:", os.RemoveAll("output"))
}

func scraping(url string) {
	c := colly.NewCollector(
		// Zenn 以外のアクセスを許可しない
		//colly.AllowedDomains(targetDomain),
		// ./cache でレスポンスをキャッシュする
		colly.CacheDir("./cache"),
		// アクセスするページの再帰の深さを設定
		//colly.MaxDepth(2),
		// ユーザーエージェントを設定
		colly.UserAgent("Sample-Scraper"),
	)

	// リクエスト間で1~2秒の時間を空ける
	c.Limit(&colly.LimitRule{
		DomainGlob:  targetDomain,
		Delay:       time.Second,
		RandomDelay: time.Second,
	})

	// エラー発生時に実行される関数
	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		prompt := []genai.Part{
			genai.Text(strings.Join(strings.Fields(e.Text), " ")),
		}
		fmt.Println(strings.Join(strings.Fields(e.Text), " "))
		geminiChat(prompt)
	})
	c.OnHTML("img", func(e *colly.HTMLElement) {
		imageURL := e.Attr("src")
		fmt.Println("\n************************************")
		//if imageURL[len(imageURL)-4:] == ".png" {
		makeImageFile(imageURL)
		fmt.Println("Image URL:", imageURL)
		if !strings.Contains(imageURL, "placeholder") {
			outputFile = "./output/" + imageURL[len(imageURL)-7:] + ".txt"
			ocr()
		}
		//} else {
		//	fmt.Println("not png", imageURL)
		//}
	})

	c.OnHTML("a", func(e *colly.HTMLElement) {
		menuURL := e.Attr("href")
		if strings.Contains(menuURL, "menu") && menuURL != url {
			fmt.Println("\n\n", menuURL, "\n")
			scraping(menuURL)
		}
	})

	c.Visit(url)
}

func ocr() {
	// Create a new Computer Vision clien
	client := computervision.New(endpoint)
	client.Authorizer = autorest.NewCognitiveServicesAuthorizer(subscription)

	// Open the image file
	file, err := os.Open(imageFilePath)
	if err != nil {
		log.Fatalf("Error opening image file: %v", err)
	}
	defer file.Close()

	fileinfo, staterr := file.Stat()

	if staterr != nil {
		fmt.Println("staterr", staterr)
		return
	}
	fmt.Println("filesize:", fileinfo.Size())
	if fileinfo.Size() < 100000 {
		fmt.Println("filesize is so small")
		return
	}

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
		log.Fatalf("Error get url", err)
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

	var prompt []genai.Part
	for _, region := range *result.Regions {
		for _, line := range *region.Lines {
			for _, word := range *line.Words {
				if word.Text != nil {
					_, err := file.WriteString(*word.Text)
					//fmt.Print(*word.Text)
					prompt = append(prompt, genai.Text(*word.Text))
					if err != nil {
						return err
					}
				}
			}
		}
	}
	geminiChat(prompt)

	return nil
}

func geminiChat(prompt []genai.Part) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiapikey))
	if err != nil {
		log.Fatalf("Error generate new client:%v", err)
	}
	model := client.GenerativeModel("gemini-pro")
	prompt = append(prompt, genai.Text("What is the lowest price. Answer true or false whether it is less than"))
	prompt = append(prompt, genai.Text(strconv.Itoa(maxprice)))
	prompt = append(prompt, genai.Text("Don't say anything other than True or False."))
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Fatalf("Error generate content:%v", err)
	}
	printResponse(resp)
}
func printResponse(resp *genai.GenerateContentResponse) {
	for _, candidate := range resp.Candidates {
		// Content が nil でないことを確認
		if candidate.Content != nil {
			fmt.Println(reflect.TypeOf(candidate.Content.Parts[0]))
			// Parts に含まれる各テキスト部分をループで処理
			for _, part := range candidate.Content.Parts {
				if part == genaiTrue {
					fmt.Println("true")
				}
			}
		}
	}
}
