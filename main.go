package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/tidwall/gjson"
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
				text := message.Text
				if strings.HasPrefix(text, "!IG ") {
					igAccount := strings.TrimPrefix(text, "!IG ")
					imgURL := findIgImg(igAccount)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(imgURL)).Do(); err != nil {
						log.Print(err)
					}
				} else {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(text+" OK!")).Do(); err != nil {
						log.Print(err)
					}
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

func findIgImg(t string) string {
	resp, err := http.Get("https://instagram.com/" + t + "/?__a=1")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "no response body is found!"
	}
	log.Print(string(body))
	pics := gjson.GetBytes(body, "graphql.user.edge_owner_to_timeline_media.edges.#.node.shortcode").Array()
	log.Print(pics)
	var randIndex int
	if len := len(pics); len > 0 {
		rand.Seed(time.Now().UnixNano())
		randIndex = rand.Intn(len)
	} else {
		return "no post found!"
	}
	return "https://instagram.com/" + t + "/p/" + pics[randIndex].String() + "/"
}
