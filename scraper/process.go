package main

import (
	"context"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
	"log"
	"regexp"
)

func GetSodanMessages(bot *traqwsbot.Bot) {
	sodanMessages, _, err := bot.
		API().
		MessageApi.
		GetMessages(context.Background(), "aff37b5f-0911-4255-81c3-b49985c8943f").
		Offset(13).
		Limit(20).
		Execute()
	if err != nil {
		log.Println(err)
	}

	//sampleMessages := []traq.Message{}
	//sampleMessages = append(sampleMessages, traq.Message{
	//	Id:        "id1",
	//	UserId:    "u1",
	//	ChannelId: "c1",
	//	Content:   "sample message",
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
		newSodan := Wiki{
			Name:        "sodan",
			Type:        "sodan",
			Content:     m.Content,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			OwnerTraqID: m.UserId,
		}
		result, err := db.Exec("INSERT INTO wikis (name, type, content, created_at, updated_at, owner_traq_id) VALUES (?, ?, ?, ?, ?, ?)",
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

		AddMessageToDB(m, int(wikiId))
	}

	GetSodanSubMessages(bot, "98ea48da-64e8-4f69-9d0d-80690b682670", 7, 52)
	GetSodanSubMessages(bot, "30c30aa5-c380-4324-b227-0ca85c34801c", 0, 32)
	GetSodanSubMessages(bot, "7ec94f1d-1920-4e15-bfc5-049c9a289692", 5, 18)
	GetSodanSubMessages(bot, "c67abb48-3fb0-4486-98ad-4b6947998ad5", 0, 21)
	GetSodanSubMessages(bot, "eb5a0035-a340-4cf6-a9e0-94ddfabe9337", 0, 2)
}

func GetSodanSubMessages(bot *traqwsbot.Bot, channelId string, offset int, limit int) {
	rsodanChannelId := "aff37b5f-0911-4255-81c3-b49985c8943f"

	messages, _, err := bot.
		API().
		MessageApi.
		GetMessages(context.Background(), channelId).
		Offset(int32(offset)).
		Limit(int32(limit)).
		Execute()
	if err != nil {
		log.Println(err)
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
			citedMessage, _, err := bot.API().MessageApi.GetMessage(context.Background(), citedMessageId).Execute()
			if err != nil {
				log.Println("failed to get cited message")
				log.Println(err)
			}
			if citedMessage.ChannelId == rsodanChannelId {
				wikiId = GetWikiIDByMessageId(citedMessageId)
			}
		}

		AddMessageToDB(m, wikiId)
	}
}

func GetWikiIDByMessageId(messageId string) int {
	var wikiId int
	err := db.Get(&wikiId, "SELECT wiki_id FROM messages WHERE message_id = ?", messageId)
	if err != nil {
		log.Println(err)
	}
	return wikiId
}

func AddMessageToDB(m traq.Message, wikiId int) {
	newMessage := Message{
		WikiID:     wikiId,
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		UserTraqID: m.UserId,
		ChannelID:  m.ChannelId,
		MessageID:  m.Id,
	}
	result, err := db.Exec("INSERT INTO messages (wiki_id, content, created_at, updated_at, user_traq_id, channel_id, message_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		newMessage.WikiID, newMessage.Content, newMessage.CreatedAt, newMessage.UpdatedAt, newMessage.UserTraqID, newMessage.ChannelID, newMessage.MessageID)
	if err != nil {
		log.Println(err)
	}
	messageId, err := result.LastInsertId()
	if err != nil {
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
		newStamp := Stamp{
			MessageID:   int(messageId),
			StampTraqID: stampId,
			Count:       count,
		}
		_, err := db.Exec("INSERT INTO messageStamps (message_id, stamp_traq_id, count) VALUES (?, ?, ?)",
			newStamp.MessageID, newStamp.StampTraqID, newStamp.Count)
		if err != nil {
			log.Println(err)
		}
	}
}
