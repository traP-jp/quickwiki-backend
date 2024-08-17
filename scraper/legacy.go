package main

import (
	"context"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"log"
)

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

func GetBotMessages(bot *traqwsbot.Bot) {
	messages, _, err := bot.
		API().
		MessageApi.
		GetMessages(context.Background(), "98ea48da-64e8-4f69-9d0d-80690b682670").
		Limit(20).
		Execute()
	if err != nil {
		log.Println(err)
	}

	for _, m := range messages {
		log.Println(m)
	}
}
