package main

import (
	"context"
	"log"
	"os"
	"strings"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
	payload "github.com/traPtitech/traq-ws-bot/payload"
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
		// _, _, err := bot.API().
		// 	MessageApi.
		// 	PostMessage(context.Background(), p.Message.ChannelID).
		// 	PostMessageRequest(traq.PostMessageRequest{
		// 		Content: "Hello",
		// 	}).
		// 	Execute()
		GetChannelIDByName(bot)
		if err != nil {
			log.Println(err)
		}
	})

	err = bot.Start()
	if err != nil {
		panic(err)
	}
}

func GetChannelIDByName(bot *traqwsbot.Bot) {
	channels, _, err := bot.API().ChannelApi.GetChannels(context.Background()).IncludeDm(false).Execute()
	if err != nil {
		panic(err)
	}

	for _, channel := range channels.Public {
		if strings.Contains(channel.Name, "random/sodan") {
			log.Printf("channleID: %s\nchannelName: %s\n", channel.Id, channel.Name)
			log.Println(channel.Children)
		}
	}
}
