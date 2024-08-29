package scraper

import (
	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traq-ws-bot/payload"
	"log"
)

func (s *Scraper) SodanMessageCreated(p *payload.MessageCreated) {
	newSodan := traq.Message{
		Id:        p.Message.ID,
		UserId:    p.Message.User.ID,
		ChannelId: p.Message.ChannelID,
		Content:   p.Message.Text,
		CreatedAt: p.Message.CreatedAt,
		UpdatedAt: p.Message.UpdatedAt,
		Stamps:    []traq.MessageStamp{},
	}
	result, err := s.db.Exec("INSERT INTO wikis (name, type, content, created_at, updated_at, owner_traq_id) VALUES (?, ?, ?, ?, ?, ?)",
		newSodan.Content, "sodan", newSodan.Content, newSodan.CreatedAt, newSodan.UpdatedAt, s.usersMap[newSodan.UserId].Name)
	if err != nil {
		log.Println("failed to insert wiki")
		log.Println(err)
	}
	wikiId, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last insert id")
		log.Println(err)
	}

	s.AddMessageToDB(newSodan, int(wikiId))
}

func SodanSubMessageCreated(p *payload.MessageCreated) {

}
