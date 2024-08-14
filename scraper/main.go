package main

import (
	"context"
	"log"
	"os"

	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	// bot.OnMessageCreated(func(p *payload.MessageCreated) {
	// 	log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
	// 	// _, _, err := bot.API().
	// 	// 	MessageApi.
	// 	// 	PostMessage(context.Background(), p.Message.ChannelID).
	// 	// 	PostMessageRequest(traq.PostMessageRequest{
	// 	// 		Content: "Hello",
	// 	// 	}).
	// 	// 	Execute()
	// 	GetMessages(p, bot)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// })
	GetChannels(bot)
}

func GetChannels(bot *traqwsbot.Bot) {
	channelID := "aff37b5f-0911-4255-81c3-b49985c8943f"
	channel, _, err := bot.API().ChannelApi.GetChannel(context.Background(), channelID).Execute()
	if err != nil {
		log.Println(err)
	}
	log.Println(channel)
	for _, c := range channel.Children {
		ch, _, err := bot.API().ChannelApi.GetChannel(context.Background(), c).Execute()
		if err != nil {
			log.Println(err)
		}
		log.Println(ch)
	}
}
