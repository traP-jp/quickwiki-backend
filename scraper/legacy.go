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
