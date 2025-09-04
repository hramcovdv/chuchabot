package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var (
	apiToken string
	sendUrl  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(".env file found")
	}

	apiToken = os.Getenv("TELEGRAM_BOT_API_TOKEN")
	if apiToken == "" {
		log.Fatal("telegram API token is required")
	}

	sendUrl = os.Getenv("TELEGRAM_BOT_SEND_URL")
	if sendUrl == "" {
		sendUrl = "http://localhost"
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Fatalf("telegram error: %s", err.Error())
	}

	// bot.Debug = true
	// log.Printf("authorized on account: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message != nil {
			err := sendJson(sendUrl, update.Message)
			if err != nil {
				log.Printf("send error: %s", err.Error())
			}
		}
	}
}

func sendJson(url string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	requestBody := bytes.NewBuffer(jsonData)
	resp, err := http.Post(url, "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
