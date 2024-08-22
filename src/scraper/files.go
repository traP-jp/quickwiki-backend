package scraper

import (
	"context"
	"log"
	"net/http"
)

func (s *Scraper) GetFile(fileId string) (*http.Response, error) {
	_, resp, err := s.bot.API().FileApi.GetFile(context.Background(), fileId).Execute()
	if err != nil {
		log.Printf("failed to get file from traq: %v", err)
		return nil, err
	}
	return resp, err
}
