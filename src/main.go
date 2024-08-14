package main

import (
	"fmt"
	"log"
	"quickwiki-backend/handler"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

func main() {
	fmt.Printf("Hello? world????\n")

	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}
	conf := mysql.Config{
		User:                 "DB_USER",
		Passwd:               "DB_PASSWORD",
		Net:                  "tcp",
		Addr:                 "DB_HOSTNAME" + ":" + "DB_PORT",
		DBName:               "DB_NAME",
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	h := handler.NewHandler(db)
	e := echo.New()

	e.GET("/ping", h.PingHandler)

	e.GET("/lecture/byFolder/id/:folderId", h.GetLectureByFolderIDHandler)
	e.GET("/lecture/byFolder/path", h.GetLectureByFolderPathHandler)
	e.GET("/lecture/folder/:folderId", h.GetLectureChildFolderHandler)
	e.GET("/lecture/lectureId", h.GetLectureHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
