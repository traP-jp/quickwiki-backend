package handler

import (
	"io"
	"log"
	"net/http"
	"quickwiki-backend/scraper"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type Handler struct {
	db      *sqlx.DB
	scraper *scraper.Scraper
}

func NewHandler(db *sqlx.DB, s *scraper.Scraper) *Handler {
	return &Handler{db: db, scraper: s}
}

// /ping
func (h *Handler) PingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

// /files/:fileId
func (h *Handler) GetFileHandler(c echo.Context) error {
	fileID := c.Param("fileId")

	resp, err := h.scraper.GetFile(fileID)
	if err != nil {
		log.Printf("failed to get file: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := c.Response()
	response.Header().Set(echo.HeaderContentType, resp.Header.Get(echo.HeaderContentType))
	response.Header().Set(echo.HeaderAccessControlExposeHeaders, "Content-Disposition")
	response.Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileID)
	response.WriteHeader(http.StatusOK)
	io.Copy(response.Writer, resp.Body)
	return c.NoContent(http.StatusOK)
}

// /stamps/:stampId
func (h *Handler) GetStampHandler(c echo.Context) error {
	stampID := c.Param("stampId")

	resp, err := h.scraper.GetStamp(stampID)
	if err != nil {
		log.Printf("failed to get stamp: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := c.Response()
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	response.Header().Set(echo.HeaderAccessControlExposeHeaders, "Content-Disposition")
	response.Header().Set(echo.HeaderContentDisposition, "attachment; filename="+stampID)
	response.WriteHeader(http.StatusOK)
	io.Copy(response.Writer, resp.Body)
	return c.NoContent(http.StatusOK)
}

// はじめのn文字を返す関数
func firstTenChars(s string, n int) string {
	// 文字列の長さがn文字未満の場合、そのまま返す
	if utf8.RuneCountInString(s) <= n {
		return s
	}

	// n文字分のルーン（文字）をスライスする
	r := []rune(s)
	return string(r[:n])
}

// mapを使用して積集合を求める関数
func intersectUsingMap(set1, set2 []int) []int {
	// マップを使ってset1の要素を記録
	setMap := make(map[int]bool)
	for _, v := range set1 {
		setMap[v] = true
	}

	// 共通部分を格納するスライス
	var intersection []int
	for _, v := range set2 {
		if setMap[v] {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

// mapを使用して和集合を求める関数
func unionUsingMap(set1, set2 []int) []int {
	// マップを使って要素を一意に保持
	setMap := make(map[int]bool)
	var union []int

	// set1の要素をマップに追加
	for _, v := range set1 {
		if !setMap[v] {
			union = append(union, v)
			setMap[v] = true
		}
	}

	// set2の要素をマップに追加
	for _, v := range set2 {
		if !setMap[v] {
			union = append(union, v)
			setMap[v] = true
		}
	}

	return union
}

func (h *Handler) SettingAllHandler(c echo.Context) error {
	h.scraper.SettingAll()
	return c.NoContent(http.StatusOK)
}

func (h *Handler) SetIndexingHandler(c echo.Context) error {
	h.scraper.SetIndexing()
	return c.NoContent(http.StatusOK)
}

type ScrapeRequest struct {
	MainChannelID string   `json:"main"`
	SubChannelIDs []string `json:"sub"`
}

func (h *Handler) ScrapingHandler(c echo.Context) error {
	req := ScrapeRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	h.scraper.GetSodanMessages(req.MainChannelID, req.SubChannelIDs)
	log.Println("finish scraping")
	return c.NoContent(http.StatusOK)
}
