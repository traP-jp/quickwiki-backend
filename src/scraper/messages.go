package scraper

import (
	"context"
	"log"
	"quickwiki-backend/model"
	"regexp"

	"github.com/traPtitech/go-traq"
)

func (s *Scraper) GetSodanMessages() {
	sodanMessages, _, err := s.bot.
		API().
		MessageApi.
		GetMessages(context.Background(), "aff37b5f-0911-4255-81c3-b49985c8943f").
		Offset(int32(16)).
		Limit(int32(20)).
		Execute()
	if err != nil {
		log.Println(err)
	}

	log.Println("--------------------")
	for _, m := range sodanMessages {
		log.Println(m.CreatedAt)
	}
	//sampleMessages := []traq.model.SodanContent_fromDB{}
	//sampleMessages = append(sampleMessages, traq.model.SodanContent_fromDB{
	//	Id:        "id1",
	//	UserId:    "u1",
	//	ChannelId: "c1",
	//	MessageContent:   "sample message",
	//	CreatedAt: time.Now(),
	//	UpdatedAt: time.Now(),
	//	Pinned:    false,
	//	Stamps: []traq.MessageStamp{
	//		{"u1", "s1", 2, time.Now(), time.Now()},
	//		{"u2", "s1", 4, time.Now(), time.Now()},
	//		{"u3", "s2", 12, time.Now(), time.Now()},
	//	},
	//})

	for _, m := range sodanMessages {
		newSodan := model.WikiContent_fromDB{
			Name:        m.Content,
			Type:        "sodan",
			Content:     m.Content,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			OwnerTraqID: s.usersMap[m.UserId].Name,
		}
		result, err := s.db.Exec("INSERT INTO wikis (name, type, content, created_at, updated_at, owner_traq_id) VALUES (?, ?, ?, ?, ?, ?)",
			newSodan.Name, newSodan.Type, newSodan.Content, newSodan.CreatedAt, newSodan.UpdatedAt, newSodan.OwnerTraqID)
		if err != nil {
			log.Println("failed to insert wiki")
			log.Println(err)
		}
		wikiId, err := result.LastInsertId()
		if err != nil {
			log.Println("failed to get last insert id")
			log.Println(err)
		}

		s.AddMessageToDB(m, int(wikiId))
	}

	s.GetSodanSubMessages("98ea48da-64e8-4f69-9d0d-80690b682670", 11, 52)
	s.GetSodanSubMessages("30c30aa5-c380-4324-b227-0ca85c34801c", 22, 32)
	s.GetSodanSubMessages("7ec94f1d-1920-4e15-bfc5-049c9a289692", 5, 18)
	s.GetSodanSubMessages("c67abb48-3fb0-4486-98ad-4b6947998ad5", 0, 21)
	s.GetSodanSubMessages("eb5a0035-a340-4cf6-a9e0-94ddfabe9337", 0, 2)

	s.updateWikisContent()
	s.setSodanTags()
	s.setIndexing()
}

func (s *Scraper) GetSodanSubMessages(channelId string, offset int, limit int) {
	rsodanChannelId := "aff37b5f-0911-4255-81c3-b49985c8943f"

	messages, _, err := s.bot.
		API().
		MessageApi.
		GetMessages(context.Background(), channelId).
		Offset(int32(offset)).
		Limit(int32(limit)).
		Execute()
	if err != nil {
		log.Println(err)
	}

	log.Println("--------------------")
	for _, m := range messages {
		log.Println(m.CreatedAt)
	}
	// reverse messages slice
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	var wikiId int
	urlOffset := len("https://q.trap.jp/messages/")
	for _, m := range messages {
		re := regexp.MustCompile(`https://q.trap.jp/messages/([^!*]{36})`)
		cites := re.FindAllString(m.Content, -1)
		if len(cites) > 0 {
			citedMessageId := cites[0][urlOffset:]
			log.Println(citedMessageId)
			citedMessage, _, err := s.bot.API().MessageApi.GetMessage(context.Background(), citedMessageId).Execute()
			if err != nil {
				log.Println("failed to get cited message")
				log.Println(err)
			}
			if citedMessage.ChannelId == rsodanChannelId {
				wikiId = s.GetWikiIDByMessageId(citedMessageId)
			}
		}

		s.AddMessageToDB(m, wikiId)
	}
}

func (s *Scraper) GetWikiIDByMessageId(messageId string) int {
	var wikiId int
	err := s.db.Get(&wikiId, "SELECT wiki_id FROM messages WHERE message_traq_id = ?", messageId)
	if err != nil {
		log.Println("failed to get wiki id")
		log.Println(err)
	}
	return wikiId
}

func (s *Scraper) AddMessageToDB(m traq.Message, wikiId int) {
	newMessage := model.SodanContent_fromDB{
		WikiID:         wikiId,
		MessageContent: m.Content,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		UserTraqID:     s.usersMap[m.UserId].Name,
		ChannelID:      m.ChannelId,
		MessageID:      m.Id,
	}
	result, err := s.db.Exec("INSERT INTO messages (wiki_id, content, created_at, updated_at, user_traq_id, channel_id, message_traq_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newMessage.WikiID, newMessage.MessageContent, newMessage.CreatedAt, newMessage.UpdatedAt, newMessage.UserTraqID, newMessage.ChannelID, newMessage.MessageID)
	if err != nil {
		log.Println("failed to insert message")
		log.Printf("%+v\nerr:%+v\n", newMessage.ID, newMessage, err)
	}
	messageId, err := result.LastInsertId()
	if err != nil {
		log.Println("failed to get last insert id")
		log.Println(err)
	}
	newMessage.ID = int(messageId)
	s.extractCitedMessage(newMessage)

	stampCount := make(map[string]int)
	for _, s := range m.Stamps {
		if _, ok := stampCount[s.StampId]; !ok {
			stampCount[s.StampId] = 0
		}
		stampCount[s.StampId] += int(s.Count)
	}

	for stampId, count := range stampCount {
		newStamp := model.Stamp_fromDB{
			MessageID:   int(messageId),
			StampTraqID: stampId,
			StampCount:  count,
		}
		_, err := s.db.Exec("INSERT INTO messageStamps (message_id, stamp_traq_id, count) VALUES (?, ?, ?)",
			newStamp.MessageID, newStamp.StampTraqID, newStamp.StampCount)
		if err != nil {
			log.Println(err)
		}
	}
}

func (s *Scraper) updateWikisContent() {
	var wikis []model.WikiContent_fromDB
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	for _, wiki := range wikis {
		var messages []model.SodanContent_fromDB
		err = s.db.Select(&messages, "SELECT * FROM messages WHERE wiki_id = ?", wiki.ID)
		if err != nil {
			log.Println("failed to get messages")
			log.Println(err)
		}

		content := ""
		for _, m := range messages {
			content += m.MessageContent + "\n"
		}

		_, err = s.db.Exec("UPDATE wikis SET content = ? WHERE id = ?", content, wiki.ID)
		if err != nil {
			log.Println("failed to update wiki")
			log.Println(err)
		}
	}
}

func (s *Scraper) extractCitedMessage(m model.SodanContent_fromDB) {
	re := regexp.MustCompile(`https://q.trap.jp/messages/([^!*]{36})`)
	cites := re.FindAllString(m.MessageContent, -1)
	for _, cite := range cites {
		messageId := cite[len("https://q.trap.jp/messages/"):]
		resp, _, err := s.bot.API().MessageApi.GetMessage(context.Background(), messageId).Execute()
		if err != nil {
			log.Println("failed to get cited message")
			log.Println(err)
		}
		citedMessage := model.CitedMessage_fromDB{
			ParentMessageID: m.ID,
			CreatedAt:       resp.CreatedAt,
			UpdatedAt:       resp.UpdatedAt,
			UserTraqID:      s.usersMap[resp.UserId].Name,
			MessageTraqID:   resp.Id,
			ChannelID:       resp.ChannelId,
			Content:         ProcessMention(resp.Content), // リンクはそのまま
		}

		_, err = s.db.Exec("INSERT INTO citedMessages (parent_message_id, created_at, updated_at, user_traq_id, message_traq_id, channel_id, content) VALUES (?, ?, ?, ?, ?, ?, ?)",
			citedMessage.ParentMessageID, citedMessage.CreatedAt, citedMessage.UpdatedAt, citedMessage.UserTraqID, citedMessage.MessageTraqID, citedMessage.ChannelID, citedMessage.Content)
		if err != nil {
			log.Println("failed to insert cited message")
			log.Println(err)
		}
	}
}
