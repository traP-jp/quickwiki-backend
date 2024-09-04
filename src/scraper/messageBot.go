package scraper

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/traPtitech/go-traq"
)

func (s *Scraper) MessageToTraQ(message string, postChanellId string) (*traq.Message, error) {
	postChanellId = "01913f8b-8c05-76a2-b51f-bb83e9e93615" //DEV_MODE
	resp, _, err := s.bot.API().
		MessageApi.
		PostMessage(context.Background(), postChanellId).
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

func (s *Scraper) MessageToDM(message string, postUserTraqId string, isUUID bool) (*traq.Message, error) {
	ret := traq.NewMessageWithDefaults()
	var userUUID string
	if isUUID {
		userUUID = postUserTraqId
	} else {
		userUUID = s.userUUIDMap[postUserTraqId]
	}
	resp, _, err := s.bot.API().MessageApi.
		PostDirectMessage(context.Background(), userUUID).
		PostMessageRequest(traq.PostMessageRequest{
			Content: message,
		}).
		Execute()
	if err != nil {
		log.Println("Failed to PostDirectMessage : ", err)
		return ret, echo.NewHTTPError(http.StatusInternalServerError, "Failed to post DM.")
	}
	return resp, echo.NewHTTPError(http.StatusOK, "successfully sent DM.")
}
