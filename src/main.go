package main

import (
	"fmt"
	"log"
	"os"
	"quickwiki-backend/handler"
	"quickwiki-backend/scraper"
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
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:               os.Getenv("DB_NAME"),
		ParseTime:            true,
		Collation:            "utf8mb4_unicode_ci",
		Loc:                  jst,
		AllowNativePasswords: true,
	}

	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	s := scraper.NewScraper(db)
	s.Scrape()

	h := handler.NewHandler(db, s)
	e := echo.New()

	e.GET("/ping", h.PingHandler)

	e.GET("/lecture/byFolder/id/:folderId", h.GetLectureByFolderIDHandler)
	e.GET("/lecture/byFolder/path", h.GetLectureByFolderPathHandler)
	e.GET("/lecture/folder/:folderId", h.GetLectureChildFolderHandler)
	e.GET("/lecture/lectureId", h.GetLectureHandler)
	e.GET("/sodan", h.GetSodanHandler)
	e.GET("/memo", h.GetMemoHandler)
	e.POST("/memo", h.PostMemoHandler)
	e.PATCH("/memo", h.PatchMemoHandler)
	e.DELETE("/memo", h.DeleteMemoHandler)
	e.POST("/wiki/tag", h.PostTagHandler)
	e.GET("/me", h.GetMeHandler)
	e.POST("/lecture", h.PostLectureHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
