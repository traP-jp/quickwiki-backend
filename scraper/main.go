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
	db       *sqlx.DB
	usersMap = make(map[string]traq.User)
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

	// get users
	users, resp, err := bot.API().UserApi.GetUsers(context.Background()).Execute()
	if err != nil {
		log.Println("failed to get users")
		log.Printf("response: %+v", resp)
		log.Fatal(err)
	}
	for _, u := range users {
		usersMap[u.Id] = u
	}

	GetSodanMessages(bot)
}
