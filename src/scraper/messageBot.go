package scraper

import (
	"context"

	"github.com/traPtitech/go-traq"
)

func (s *Scraper) MessageToTraQ(message string) error {
	botChannelId := "01913f8b-8c05-76a2-b51f-bb83e9e93615"
	_, _, err := s.bot.API().
		MessageApi.
		PostMessage(context.Background(), botChannelId).
		PostMessageRequest(traq.PostMessageRequest{
			Content: message,
		}).
		Execute()
	return err
}
