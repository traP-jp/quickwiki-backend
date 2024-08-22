package scraper

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"log"
	"os"
)

type Scraper struct {
	db       *sqlx.DB
	usersMap map[string]traq.User
	bot      *traqwsbot.Bot
}

func NewScraper(db *sqlx.DB) *Scraper {
	return &Scraper{
		db:       db,
		usersMap: make(map[string]traq.User),
	}
}

func (s *Scraper) Scrape() {
	// setting bot
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}
	log.Println("bot connected")
	s.bot = bot

	// get users
	users, resp, err := bot.API().UserApi.GetUsers(context.Background()).Execute()
	if err != nil {
		log.Println("failed to get users")
		log.Printf("response: %+v", resp)
		log.Fatal(err)
	}
	for _, u := range users {
		s.usersMap[u.Id] = u
	}

	s.GetSodanMessages()
}
