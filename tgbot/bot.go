package tgbot

import (
	"context"
	"deluge/config"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/reply"
)

var demoReplyKeyboard *reply.ReplyKeyboard
var doc models.Document
var cfg *config.Config

func New(config *config.Config) (*bot.Bot, error) {
	cfg = config
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(cfg.TelegramToken, opts...)
	if err != nil {
		panic(err)
	}
	return b, nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	document := update.Message.Document
	if document != nil {
		name := document.FileName
		if strings.Contains(name, ".torrent") {
			doc = *document
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        "Какой это тип?",
				ReplyMarkup: demoReplyKeyboard,
			})
		}
	}

	println("got file:", doc.FileName)
}

func InitReplyKeyboard(b *bot.Bot) {
	demoReplyKeyboard = reply.New(
		reply.WithPrefix("reply_keyboard"),
		reply.IsSelective(),
		reply.IsOneTimeKeyboard(),
	).
		Button("Фильм", b, bot.MatchTypeExact, onReplyKeyboardSelect).
		Row().
		Button("Сериал", b, bot.MatchTypeExact, onReplyKeyboardSelect)
}

func onReplyKeyboardSelect(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}
	switch update.Message.Text {
	case "Фильм":
		err := downloadFile(ctx, b, &doc, cfg.FilmsPath+doc.FileName)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Error:\n" + err.Error(),
			})
			println(err.Error())
		}
	case "Сериал":
		err := downloadFile(ctx, b, &doc, cfg.ShowsPath+doc.FileName)
		if err != nil {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Error:\n" + err.Error(),
			})
			println(err.Error())
		}

	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        "Какой это тип?",
			ReplyMarkup: demoReplyKeyboard,
		})

	}
}

func downloadFile(ctx context.Context, b *bot.Bot, doc *models.Document, path string) error {
	file, err := b.GetFile(ctx, &bot.GetFileParams{FileID: doc.FileID})
	url := b.FileDownloadLink(file)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
