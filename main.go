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

type Core struct {
	ch chan tgbotapi.Update
}

func NewCore(size int) *Core {
	return &Core{
		ch: make(chan tgbotapi.Update, size),
	}
}

func (c *Core) update(u tgbotapi.Update) {
	c.ch <- u
}

func (c *Core) loop() {
	for u := range c.ch {
		if !u.Message.IsCommand() {
			err := sendJson(sendUrl, u.Message)
			if err != nil {
				log.Printf("error: %s", err.Error())
			}
		}
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Print(".env file not found")
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

	core := NewCore(10)
	go core.loop()

	for update := range updates {
		if update.Message != nil {
			if !update.Message.IsCommand() {
				core.update(update)
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
