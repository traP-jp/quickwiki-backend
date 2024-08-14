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
		Addr:                 "mariadb:3306",
		DBName:               "quickwiki",
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	db, err = sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

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
	GetChannels(bot)
}

func GetChannels(bot *traqwsbot.Bot) {
	channelID := "aff37b5f-0911-4255-81c3-b49985c8943f"
	channel, _, err := bot.API().ChannelApi.GetChannel(context.Background(), channelID).Execute()
	if err != nil {
		log.Println(err)
	}
	log.Println(channel)
	for _, c := range channel.Children {
		ch, _, err := bot.API().ChannelApi.GetChannel(context.Background(), c).Execute()
		if err != nil {
			log.Println(err)
		}
		log.Println(ch)
	}
}

func GetSodanMessages(bot *traqwsbot.Bot) {
	sodanMessages, _, err := bot.API().MessageApi.GetMessages(context.Background(), "aff37b5f-0911-4255-81c3-b49985c8943f").Offset(13).Limit(20).Order("asc").Execute()
	if err != nil {
		log.Println(err)
	}
	
	

	for _, m := range sodanMessages {
		log.Println(m)
		newSodan := &Wiki{
			Name:        "sodan",
			Type:        "sodan",
			Content:     m.Content,
			CreatedAt:   m.CreatedAt,
			UpdatedAt:   m.UpdatedAt,
			OwnerTraqID: m.UserId,
		}
		result, err := db.Exec("INSERT INTO wikis (name, type, content, created_at, updated_at, owner_traq_id) VALUES (?, ?, ?, ?, ?, ?)", newSodan.Name, newSodan.Type, newSodan.Content, newSodan.CreatedAt, newSodan.UpdatedAt, newSodan.OwnerTraqID)
		if err != nil {
			log.Println(err)
		}
		wikiId := result.LastInsertId()

		AddMessageToDB(m, wikiId)
	}
}

func AddMessageToDB(m traq.Message, wikiId int) {
	newMessage := &Message{
		WikiID:     wikiId,
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
	messageId := result.LastInsertId()

	stampCount := make(map[string]int)
	for _, s := range m.Stamps {
		if _, ok := stampCount[s.StampId]; !ok {
			stampCount[s.StampId] = 0
		}
		stampCount[s.StampId]++
	}

	for stampId, count := range stampCount {
		newStamp := &Stamp{
			MessageID:   messageId,
			StampTraqID: stampId,
			Count:       count,
		}
		_, err := db.Exec("INSERT INTO stamps (message_id, stamp_traq_id, count) VALUES (?, ?, ?)", newStamp.MessageID, newStamp.StampTraqID, newStamp.Count)
		if err != nil {
			log.Println(err)
		}
	}
}