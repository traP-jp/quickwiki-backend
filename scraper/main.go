package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/traPtitech/go-traq"
	traqwsbot "github.com/traPtitech/traq-ws-bot"
)

var (
	db *sqlx.DB
)

func main() {
	// setting bot
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("TRAQ_BOT_TOKEN"),
	})
	if err != nil {
		panic(err)
	}

	// setting db
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	conf := mysql.Config{
		User:                 "root",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "quickwiki",
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	db, err = sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Println("failed to open db")
		log.Fatal(err)
	}
	log.Println("connected")

	// bot.OnMessageCreated(func(p *payload.MessageCreated) {
	// 	log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
	// 	// _, _, err := bot.API().
	// 	// 	MessageApi.
	// 	// 	PostMessage(context.Background(), p.Message.ChannelID).
	// 	// 	PostMessageRequest(traq.PostMessageRequest{
	// 	// 		Content: "Hello",
	// 	// 	}).
	// 	// 	Execute()
	// 	GetMessages(p, bot)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// })
	//GetSodanMessages(bot)
	GetBotMessages(bot)
}

func GetBotMessages(bot *traqwsbot.Bot) {
	messages, _, err := bot.
		API().
		MessageApi.
		GetMessages(context.Background(), "98ea48da-64e8-4f69-9d0d-80690b682670").
		Limit(20).
		Execute()
	if err != nil {
		log.Println(err)
	}

	for _, m := range messages {
		log.Println(m)
	}
}

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
		result, err := db.Exec("INSERT INTO wikis (name, type, content, created_at, updated_at, owner_traq_id) VALUES (?, ?, ?, ?, ?, ?)", newSodan.Name, newSodan.Type, newSodan.Content, newSodan.CreatedAt, newSodan.UpdatedAt, newSodan.OwnerTraqID)
		if err != nil {
			log.Println("failed to insert wiki")
			log.Println(err)
		}
		wikiId, err := result.LastInsertId()
		if err != nil {
			log.Println("failed to get last insert id")
			log.Println(err)
		}

		AddMessageToDB(m, wikiId)
	}
}

func AddMessageToDB(m traq.Message, wikiId int64) {
	newMessage := Message{
		WikiID:     int(wikiId),
		Content:    m.Content,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		UserTraqID: m.UserId,
		ChannelID:  m.ChannelId,
	}
	result, err := db.Exec("INSERT INTO messages (wiki_id, content, created_at, updated_at, user_traq_id, channel_id) VALUES (?, ?, ?, ?, ?, ?)", newMessage.WikiID, newMessage.Content, newMessage.CreatedAt, newMessage.UpdatedAt, newMessage.UserTraqID, newMessage.ChannelID)
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
		_, err := db.Exec("INSERT INTO messageStamps (message_id, stamp_traq_id, count) VALUES (?, ?, ?)", newStamp.MessageID, newStamp.StampTraqID, newStamp.Count)
		if err != nil {
			log.Println(err)
		}
	}
}
