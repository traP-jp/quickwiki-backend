package scraper

import (
	"context"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

type Scraper struct {
	db                *sqlx.DB
	usersMap          map[string]traq.User
	userUUIDMap       map[string]string
	usersDisplayNames map[string]string
	bot               *traqwsbot.Bot
}

func NewScraper(db *sqlx.DB) *Scraper {
	return &Scraper{
		db:                db,
		usersMap:          make(map[string]traq.User),
		userUUIDMap:       make(map[string]string),
		usersDisplayNames: make(map[string]string),
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
	users, resp, err := bot.API().UserApi.GetUsers(context.Background()).IncludeSuspended(true).Execute()
	if err != nil {
		log.Println("failed to get users")
		log.Printf("response: %+v", resp)
		log.Fatal(err)
	}
	for _, u := range users {
		s.usersMap[u.Id] = u
		s.userUUIDMap[u.Name] = u.Id
		s.usersDisplayNames[u.Name] = u.DisplayName
	}

	// s.GetSodanMessages()

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		channelId := p.Message.ChannelID
		if channelId == "aff37b5f-0911-4255-81c3-b49985c8943f" {
			s.SodanMessageCreated(p)
		} else if channelId == "98ea48da-64e8-4f69-9d0d-80690b682670" || channelId == "30c30aa5-c380-4324-b227-0ca85c34801c" || channelId == "7ec94f1d-1920-4e15-bfc5-049c9a289692" || channelId == "c67abb48-3fb0-4486-98ad-4b6947998ad5" || channelId == "eb5a0035-a340-4cf6-a9e0-94ddfabe9337" || channelId == "5b857a8d-03b5-4c25-92d9-bc01f3defe84" {
			s.SodanSubMessageCreated(p)
		}
	})
}

func (s *Scraper) StartBot() {
	err := s.bot.Start()
	if err != nil {
		panic(err)
	}
}
