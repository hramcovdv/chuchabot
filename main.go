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
	ch chan *tgbotapi.Message
}

func NewCore() *Core {
	core := &Core{
		ch: make(chan *tgbotapi.Message, 10),
	}

	go core.loop()

	return core
}

func (c *Core) loop() {
	for msg := range c.ch {
		if !msg.IsCommand() {
			err := sendJson(sendUrl, msg)
			if err != nil {
				log.Printf("error: %s", err.Error())
			}
		}
	}
}

func (c *Core) Recive(msg *tgbotapi.Message) {
	c.ch <- msg
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

	core := NewCore()

	for update := range updates {
		if update.Message != nil {
			core.Recive(update.Message)
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
