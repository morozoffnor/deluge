package server

import (
	"context"
	"deluge/config"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-telegram/bot"
)

type Event struct {
	Type string `json:"event_type"`
	Name string `json:"name"`
	Path string `json:"path"`
	Id   string `json:"id"`
}

func New(cfg *config.Config, b *bot.Bot) *http.Server {
	return &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: newHandler(cfg, b),
	}
}

func newHandler(cfg *config.Config, b *bot.Bot) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		event := Event{}
		err = json.Unmarshal(buf, &event)
		if err != nil {
			println(err)
			return
		}
		sendNotification(cfg, b, &event)
		println(string(buf))
	})
}

func sendNotification(cfg *config.Config, b *bot.Bot, e *Event) {
	var eventType string
	switch e.Type {
	case "added":
		eventType = "Torrent added"
	case "completed":
		eventType = "Torrent completed"
	default:
		eventType = "Unknown event"
	}
	msg, err := b.SendMessage(context.Background(), &bot.SendMessageParams{
		ChatID:    cfg.TelegramChatID,
		Text:      bot.EscapeMarkdownUnescaped(e.Name + "\n\n" + eventType),
		ParseMode: "MarkdownV2",
	})
	println(msg, err)
}
