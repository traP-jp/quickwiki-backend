package scraper

import (
	"context"
	"log"
	"net/http"
)

func (s *Scraper) GetFile(fileId string) (*http.Response, error) {
	_, resp, err := s.bot.API().FileApi.GetFile(context.Background(), fileId).Execute()
	if err != nil {
		log.Printf("failed to get file '%s' from traq: %v", fileId, err)
		return nil, err
	}
	return resp, err
}

func (s *Scraper) GetStamp(stampId string) (*http.Response, error) {
	stamp, resp, err := s.bot.API().StampApi.GetStamp(context.Background(), stampId).Execute()
	if err != nil {
		log.Printf("failed to get stamp from traq: %v", err)
		log.Printf("response: %v", resp)
		return nil, err
	}

	resp, err = s.GetFile(stamp.FileId)
	if err != nil {
		log.Printf("failed to get stamp '%s' file: %v", stamp.Name, err)
		return nil, err
	}

	return resp, err
}
