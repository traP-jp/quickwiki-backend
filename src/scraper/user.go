package scraper

import (
	"fmt"
	"quickwiki-backend/model"
)

// return without iconuri
func (s *Scraper) GetUserDetail(userTraqID string) (model.Me_Response, error) {
	displayName, ok := s.usersDisplayNames[userTraqID]
	if !ok {
		return model.Me_Response{}, fmt.Errorf("user not found")
	}

	return model.Me_Response{
		TraqID:      userTraqID,
		DisplayName: displayName,
	}, nil
}
