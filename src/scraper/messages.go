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

func (s *Scraper) GetSodanMessages(mainChannelId string, subChannelId []string) {
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
	sodanMessages, err := s.getMessages(mainChannelId)
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

	for i, channelId := range subChannelId {
		s.GetSodanSubMessages(channelId, mainChannelId)
		log.Printf("finished scraping subchannel %d\n", i+1)
	}

	//s.MergeWikisContent()
	//s.SetSodanTags()
	//s.RemoveMentions()
	//s.SetIndexing()
}

func (s *Scraper) GetSodanSubMessages(channelId string, parentId string) {

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
			} else if citedMessage.ChannelId == parentId {
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

func (s *Scraper) UpdateMessageToDB(m traq.Message, wikiId int) {
	newMessage := model.SodanContent_fromDB{
		WikiID:         wikiId,
		MessageContent: m.Content,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		UserTraqID:     s.usersMap[m.UserId].Name,
		ChannelID:      m.ChannelId,
		MessageID:      m.Id,
	}
	_, err := s.db.Exec("UPDATE messages SET content = ?, updated_at = ? WHERE message_traq_id = ?",
		newMessage.MessageContent, newMessage.UpdatedAt, newMessage.MessageID)
	if err != nil {
		log.Println("failed to update message")
		log.Printf("%+v\nerr:%+v\n", newMessage.ID, err)
	}
	err = s.db.Get(&newMessage.ID, "SELECT id FROM messages WHERE message_traq_id = ?", newMessage.MessageID)
	if err != nil {
		log.Println("failed to get message id")
		log.Println(err)
	}

	stampCount := make(map[string]int)
	for _, s := range m.Stamps {
		if _, ok := stampCount[s.StampId]; !ok {
			stampCount[s.StampId] = 0
		}
		stampCount[s.StampId] += int(s.Count)
	}

	for stampId, count := range stampCount {
		newStamp := model.Stamp_fromDB{
			MessageID:   newMessage.ID,
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

func (s *Scraper) RemoveMentionFromMessage() {
	var messages []model.SodanContent_fromDB
	err := s.db.Select(&messages, "SELECT * FROM messages")
	if err != nil {
		log.Println("failed to get messages")
		log.Println(err)
	}

	for _, m := range messages {
		content := ProcessMention(m.MessageContent)
		_, err = s.db.Exec("UPDATE messages SET content = ? WHERE id = ?", content, m.ID)
		if err != nil {
			log.Println("failed to update message")
			log.Println(err)
		}
	}
}

func (s *Scraper) FixTitle() {
	var wikis []model.WikiContent_fromDB
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	for _, wiki := range wikis {
		name := ProcessMentionAll(ProcessLink(removeNewLine(removeCodeBlock(removeTeX(wiki.Name)))))
		r := []rune(name)
		if len(r) > 50 {
			name = string(r[:50]) + "..."
		}
		log.Println(name)
		_, err = s.db.Exec("UPDATE wikis SET name = ? WHERE id = ?", name, wiki.ID)
		if err != nil {
			log.Printf("failed to update wiki: %v", err)
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

func (s *Scraper) SettingAll() {
	s.MergeWikisContent()
	log.Println("finished merging wikis content")
	//s.SetSodanTags()
	log.Println("finished setting sodan tags")
	s.FixTitle()
	log.Println("finished fixing title")
	s.RemoveMentions()
	log.Println("finished removing mentions")
	s.RemoveMentionFromMessage()
	log.Println("finished removing mention from message")
	s.SetIndexing()
	log.Println("finished setting indexing")
}
