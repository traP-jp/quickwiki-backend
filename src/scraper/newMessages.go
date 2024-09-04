package scraper

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"quickwiki-backend/model"
	"quickwiki-backend/search"
	"regexp"

	"github.com/traPtitech/go-traq"
	"github.com/traPtitech/traq-ws-bot/payload"
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
	//SodanSubMessageCreatedかつ質問者とsub投稿者が違うとき、DMに通知する
	s.responseNotification(p, wikiId)
}
func (s *Scraper) responseNotification(p *payload.MessageCreated, wikiId int) {
	var messages model.SodanContent_fromDB
	err := s.db.Get(&messages, "SELECT * FROM messages WHERE wiki_id = ? ORDER BY created_at ASC LIMIT 1", wikiId)
	if err != nil {
		log.Println("failed to get messages")
		log.Println(err)
	}
	//anon-sodanを使用していればanonSodanの表に残っているはず,そうでなければ投稿者の判定をそのまま行う
	var anonSodan model.AnonSodans_fromDB
	var userUUID string
	err = s.db.Get(&anonSodan, "SELECT * FROM anonSodan WHERE message_traq_id = ?", messages.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			userUUID = s.userUUIDMap[messages.UserTraqID]
		}
		log.Printf("[in delete memo]failed to get wikiContent: %s\n", err)
		return
	}
	userUUID = s.userUUIDMap[anonSodan.UserTraqID]
	if p.Message.User.ID == userUUID {
		return
	}
	message := "You got a reply to your SODAN in the random/sodan channel.\nhttps://q.trap.jp/messages/" + p.Message.ID
	s.MessageToDM(message, p.Message.User.ID, true)
}

func (s *Scraper) getWikiId(channelId string) int {
	var wikiId int
	urlOffset := len("https://q.trap.jp/messages/")
	rsodanChannelId := "aff37b5f-0911-4255-81c3-b49985c8943f"

	var messages []model.SodanContent_fromDB
	err := s.db.Select(&messages, "SELECT * FROM messages WHERE channel_id = ? ORDER BY created_at DESC LIMIT 30", channelId)
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
}
