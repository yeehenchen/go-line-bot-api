package main

import (
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	log.Println("Bot:", bot)
	bot, err = linebot.New(os.Getenv("LINE_SECRET"), os.Getenv("LINE_TOKEN"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/webhook", requestHandler)
	port := getPort()
	http.ListenAndServe(port, nil)
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(404)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				_, err := bot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text+" OK!")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func getPort() string {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return ":" + port
}
