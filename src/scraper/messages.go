package scraper

import (
	"context"
	"fmt"
	"log"
	"os"
	"quickwiki-backend/model"
	"regexp"
	"time"

	"github.com/traPtitech/go-traq"
)

func (s *Scraper) getMessages(channelId string) ([]traq.Message, error) {
	var messages []traq.Message
	for i := 0; ; i++ {
		res, r, err := s.bot.API().MessageApi.GetMessages(context.Background(), channelId).Limit(int32(200)).Offset(int32(i * 200)).Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `ChannelApi.GetMessages``: %v\n", err)
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		}

		if err != nil {
			return nil, err
		}
		if len(res) == 0 {
			break
		}
		messages = append(messages, res...)
		time.Sleep(time.Millisecond * 1000)
		log.Printf("fetched %d messages\n", len(messages))
	}

	return messages, nil
}

func (s *Scraper) GetSodanMessages() {
	//sodanMessages, resp, err := s.bot.
	//	API().
	//	MessageApi.
	//	GetMessages(context.Background(), "aff37b5f-0911-4255-81c3-b49985c8943f").
	//	//Offset(int32(17)).
	//	Limit(int32(3000)).
	//	Execute()
	//if err != nil {
	//	log.Println(err)
	//	log.Println(resp)
	//}
	sodanMessages, err := s.getMessages("aff37b5f-0911-4255-81c3-b49985c8943f")
	if err != nil {
		log.Println(err)
	}

	log.Println("--------------------")
	log.Printf("sodan messages count: %d\n", len(sodanMessages))

	for i, m := range sodanMessages {
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

		if i%100 == 0 {
			log.Printf("inserted %d wikis\n", i)
		}

		s.AddMessageToDB(m, int(wikiId))
	}

	s.GetSodanSubMessages("98ea48da-64e8-4f69-9d0d-80690b682670", 40, 52)
	log.Println("sodan messages scraped")
	s.GetSodanSubMessages("98ea48da-64e8-4f69-9d0d-80690b682670", 48, 52)
	log.Println("sodan sub messages 1 scraped")
	s.GetSodanSubMessages("30c30aa5-c380-4324-b227-0ca85c34801c", 22, 32)
	log.Println("sodan sub messages 2 scraped")
	s.GetSodanSubMessages("7ec94f1d-1920-4e15-bfc5-049c9a289692", 5, 18)
	log.Println("sodan sub messages 3 scraped")
	s.GetSodanSubMessages("c67abb48-3fb0-4486-98ad-4b6947998ad5", 0, 21)
	log.Println("sodan sub messages 4 scraped")
	s.GetSodanSubMessages("eb5a0035-a340-4cf6-a9e0-94ddfabe9337", 0, 2)
	log.Println("sodan sub messages 5 scraped")
	s.GetSodanSubMessages("5b857a8d-03b5-4c25-92d9-bc01f3defe84", 0, 2)

	//s.MergeWikisContent()
	//s.SetSodanTags()
	//s.RemoveMentions()
	//s.SetIndexing()
}

func (s *Scraper) GetSodanSubMessages(channelId string, offset int, limit int) {
	rsodanChannelId := "aff37b5f-0911-4255-81c3-b49985c8943f"

	//messages, _, err := s.bot.
	//	API().
	//	MessageApi.
	//	GetMessages(context.Background(), channelId).
	//	//Offset(int32(offset)).
	//	Limit(int32(3000)).
	//	Execute()
	messages, err := s.getMessages(channelId)
	if err != nil {
		log.Println(err)
	}

	log.Println("--------------------")
	// reverse messages slice
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	wikiId := 1
	urlOffset := len("https://q.trap.jp/messages/")
	for i, m := range messages {
		re := regexp.MustCompile(`https://q.trap.jp/messages/([^!*]{36})`)
		cites := re.FindAllString(m.Content, -1)
		if len(cites) > 0 {
			citedMessageId := cites[0][urlOffset:]
			citedMessage, _, err := s.bot.API().MessageApi.GetMessage(context.Background(), citedMessageId).Execute()
			if err != nil {
				log.Println("failed to get cited message")
				log.Println(err)
			} else if citedMessage.ChannelId == rsodanChannelId {
				wikiId = s.GetWikiIDByMessageId(citedMessageId)
			}
		}

		if i%200 == 0 {
			log.Printf("inserted %d wikis\n", i)
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
		log.Printf("%+v\nerr:%+v\n", newMessage.ID, err)
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

func (s *Scraper) MergeWikisContent() {
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
			content += m.MessageContent
		}

		_, err = s.db.Exec("UPDATE wikis SET content = ? WHERE id = ?", content, wiki.ID)
		if err != nil {
			log.Println("failed to update wiki")
			log.Println(err)
		}
	}
}

func (s *Scraper) RemoveMentions() {
	var wikis []model.WikiContent_fromDB
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	for _, wiki := range wikis {
		content := ProcessMention(wiki.Content)
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
			continue
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
