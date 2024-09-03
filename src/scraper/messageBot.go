package scraper

import (
	"context"

	"github.com/traPtitech/go-traq"
)

func (s *Scraper) MessageToTraQ(message string, PostChanellId string) (*traq.Message, error) {
	PostChanellId = "01913f8b-8c05-76a2-b51f-bb83e9e93615" //DEV_MODE
	resp, _, err := s.bot.API().
		MessageApi.
		PostMessage(context.Background(), PostChanellId).
		PostMessageRequest(traq.PostMessageRequest{
			Content: message,
		}).
		Execute()
	return resp, err
}

func (s *Scraper) MessageEditOnTraQ(message string, editMessageId string) error {
	editMessageId = "0191a3c2-6659-7c4e-a03f-6a37d080cd10" //DEV_MODE
	_, err := s.bot.API().MessageApi.
		EditMessage(context.Background(), editMessageId).
		PostMessageRequest(traq.PostMessageRequest{
			Content: message,
		}).
		Execute()
	return err
}
