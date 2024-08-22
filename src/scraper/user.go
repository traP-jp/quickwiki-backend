package scraper

import (
	"fmt"
	"log"
	"quickwiki-backend/model"
)

// return without iconuri
func (s *Scraper) GetUserDetail(userTraqID string) (model.Me_Response, error) {
	user, ok := s.usersMap[userTraqID]
	if !ok {
		log.Println("user not found")
		return model.Me_Response{}, fmt.Errorf("user not found")
	}

	return model.Me_Response{
		TraqID: user.Name,
		Name:   user.DisplayName,
	}, nil
}
