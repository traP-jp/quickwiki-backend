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
		User:                 os.Getenv("NS_MARIADB_USER"),
		Passwd:               os.Getenv("NS_MARIADB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("NS_MARIADB_HOSTNAME") + ":" + os.Getenv("NS_MARIADB_PORT"),
		DBName:               os.Getenv("NS_MARIADB_DATABASE"),
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

	e.POST("/lecture", h.PostLectureHandler)
	e.GET("/lecture/byFolder/id/:folderId", h.GetLectureByFolderIDHandler)
	e.GET("/lecture/byFolder/path", h.GetLectureByFolderPathHandler)
	e.GET("/lecture/folder/:folderId", h.GetLectureChildFolderHandler)
	e.GET("/lecture/:lectureId", h.GetLectureHandler)

	e.GET("/sodan", h.GetSodanHandler)
	e.GET("/memo", h.GetMemoHandler)
	e.POST("/memo", h.PostMemoHandler)
	e.PATCH("/memo", h.PatchMemoHandler)
	e.DELETE("/memo", h.DeleteMemoHandler)
	e.GET("/tag", h.GetTagsHandler)
	e.POST("/wiki/search", h.SearchHandler)
	e.GET("/wiki/tag", h.GetWikiByTagHandler)
	e.POST("/wiki/tag", h.PostTagHandler)
	e.PATCH("/wiki/tag", h.EditTagHandler)
	e.DELETE("/wiki/tag", h.DeleteTagHandler)

	e.GET("/me", h.GetMeHandler)
	e.GET("/wiki/user", h.GetUserWikiHandelr)
	e.GET("/wiki/user/favorite", h.GetUserFavoriteWikiHandler)
	e.POST("/wiki/user/favorite", h.PostUserFavoriteWikiHandler)
	e.DELETE("/wiki/user/favorite", h.DeleteUserFavoriteWikiHandler)

	e.POST("/anon-sodan", h.PostMessageToTraQ)
	e.PATCH("/anon-sodan", h.PatchMessageToTraQ)
	e.POST("/anon-sodan/replies", h.PostRepliesToTraQ)

	e.GET("/files/:fileId", h.GetFileHandler)
	e.GET("/stamps/:stampId", h.GetStampHandler)

	e.GET("/setting/index", h.SetIndexingHandler)
	//e.POST("/setting/messages", h.ScrapingHandler)
	e.GET("/setting/all", h.SettingAllHandler)

	e.Logger.Fatal(e.Start(":8080"))
	//s.StartBot()
}
