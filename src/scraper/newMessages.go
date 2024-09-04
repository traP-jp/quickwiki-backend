package scraper

import (
	"context"
	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traq-ws-bot/payload"
	"log"
	"quickwiki-backend/model"
	"quickwiki-backend/search"
	"regexp"
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
	s.addMessageTag(int(wikiId))
	s.addMessageIndex(int(wikiId))
}

func (s *Scraper) SodanSubMessageCreated(p *payload.MessageCreated) {
	rsodanChannelId := "aff37b5f-0911-4255-81c3-b49985c8943f"

	wikiId := 0
	urlOffset := len("https://q.trap.jp/messages/")
	re := regexp.MustCompile(`https://q.trap.jp/messages/([^!*]{36})`)
	cites := re.FindAllString(p.Message.Text, -1)
	if len(cites) > 0 {
		citedMessageId := cites[0][urlOffset:]
		citedMessage, _, err := s.bot.API().MessageApi.GetMessage(context.Background(), citedMessageId).Execute()
		if err != nil {
			log.Println("failed to get cited message")
			log.Println(err)
		}
		if citedMessage.ChannelId == rsodanChannelId {
			s.registerWiki(p.Message.ChannelID)
			wikiId = s.GetWikiIDByMessageId(citedMessageId)
		}
	}

	if wikiId == 0 {
		wikiId = s.getWikiId(p.Message.ChannelID)
	}

	s.AddMessageToDB(traq.Message{
		Id:        p.Message.ID,
		UserId:    p.Message.User.ID,
		ChannelId: p.Message.ChannelID,
		Content:   p.Message.Text,
		CreatedAt: p.Message.CreatedAt,
		UpdatedAt: p.Message.UpdatedAt,
		Stamps:    []traq.MessageStamp{},
	}, wikiId)
	s.addMessageTag(wikiId)
	s.removeMentionSingle(wikiId)
	s.addMessageIndex(wikiId)
}

func (s *Scraper) getWikiId(channelId string) int {
	var wikiId int
	urlOffset := len("https://q.trap.jp/messages/")
	rsodanChannelId := "aff37b5f-0911-4255-81c3-b49985c8943f"

	var messages []model.SodanContent_fromDB
	err := s.db.Select(&messages, "SELECT * FROM messages WHERE channel_id = ? ORDER BY created_at DESC LIMIT 100", channelId)
	if err != nil {
		log.Println("failed to get messages")
		log.Println(err)
	}

	// check where starts the sodan
	for _, m := range messages {
		re := regexp.MustCompile(`https://q.trap.jp/messages/([^!*]{36})`)
		cites := re.FindAllString(m.MessageContent, -1)
		if len(cites) > 0 {
			citedMessageId := cites[0][urlOffset:]
			citedMessage, _, err := s.bot.API().MessageApi.GetMessage(context.Background(), citedMessageId).Execute()
			if err != nil {
				log.Println("failed to get cited message")
				log.Println(err)
			}
			if citedMessage.ChannelId == rsodanChannelId {
				wikiId = m.WikiID
				break
			}
		}
	}

	return wikiId
}

func (s *Scraper) registerWiki(channelId string) {
	wikiId := s.getWikiId(channelId)

	s.mergeWikiContent(wikiId)
	s.addMessageTag(wikiId)
	s.removeMentionSingle(wikiId)
	s.addMessageIndex(wikiId)
}

func (s *Scraper) mergeWikiContent(wikiId int) {
	var wiki model.WikiContent_fromDB
	err := s.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiId)
	if err != nil {
		log.Println("failed to get wiki")
		log.Println(err)
	}

	var messages []model.SodanContent_fromDB
	err = s.db.Select(&messages, "SELECT * FROM messages WHERE wiki_id = ?", wikiId)
	if err != nil {
		log.Println("failed to get messages")
		log.Println(err)
	}

	content := ""
	for _, m := range messages {
		content += m.MessageContent
	}

	_, err = s.db.Exec("UPDATE wikis SET content = ? WHERE id = ?", content, wikiId)
	if err != nil {
		log.Println("failed to update wiki")
		log.Println(err)
	}
}

func (s *Scraper) addMessageTag(wikiId int) {
	var wiki model.WikiContent_fromDB
	err := s.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiId)
	if err != nil {
		log.Println("failed to get wiki")
		log.Println(err)
	}

	s.setTag([]model.WikiContent_fromDB{wiki})
}

func (s *Scraper) addMessageIndex(wikiId int) {
	var wiki model.WikiContent_fromDB
	err := s.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiId)
	if err != nil {
		log.Println("failed to get wiki")
		log.Println(err)
	}

	indexData := []search.IndexData{
		{
			ID:             wiki.ID,
			Type:           wiki.Type,
			Title:          wiki.Name,
			OwnerTraqID:    wiki.OwnerTraqID,
			MessageContent: wiki.Content,
		},
	}

	search.Indexing(indexData)
}

func (s *Scraper) removeMentionSingle(wikiId int) {
	var wiki model.WikiContent_fromDB
	err := s.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiId)
	if err != nil {
		log.Println("failed to get wiki")
		log.Println(err)
	}

	text := ProcessMention(wiki.Content)
	_, err = s.db.Exec("UPDATE wikis SET content = ? WHERE id = ?", text, wikiId)
	if err != nil {
		log.Println("failed to update wiki")
		log.Println(err)
	}

	var messages []model.SodanContent_fromDB
	err = s.db.Select(&messages, "SELECT * FROM messages WHERE wiki_id = ?", wikiId)
	if err != nil {
		log.Println("failed to get messages")
		log.Println(err)
	}

	for _, m := range messages {
		text := ProcessMention(m.MessageContent)
		_, err = s.db.Exec("UPDATE messages SET message_content = ? WHERE id = ?", text, m.ID)
		if err != nil {
			log.Println("failed to update message")
			log.Println(err)
		}
	}
}
