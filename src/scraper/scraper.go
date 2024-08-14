package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	traq "github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	payload "github.com/traPtitech/traq-ws-bot/payload"
)

type Scraper struct {
	db *sqlx.DB
}

func NewScraper(db *sqlx.DB) *Scraper {
	return &Scraper{db: db}
}

func (s *Scraper) StartBot(c echo.Context) error {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
		_, _, err := bot.API().
			MessageApi.
			PostMessage(context.Background(), p.Message.ChannelID).
			PostMessageRequest(traq.PostMessageRequest{
				Content: "Hello",
			}).
			Execute()
		if err != nil {
			log.Println(err)
		}
	})

	err = bot.Start()
	if err != nil {
		panic(err)
	}

	return c.String(http.StatusOK, "BOT started")
}
