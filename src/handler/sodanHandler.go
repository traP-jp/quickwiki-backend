package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"strconv"

	"github.com/labstack/echo"
)

func (h *Handler) PostMessageToTraQ(c echo.Context) error {
	var message model.MessageToTraQ_POST
	err := c.Bind(&message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	channelId := "aff37b5f-0911-4255-81c3-b49985c8943f"
	err = h.scraper.MessageToTraQ(message.Content, channelId)

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, h.GetSodanHandler(c))
}

func (h *Handler) PatchMessageToTraQ(c echo.Context) error {
	wikiId, err := strconv.Atoi(c.QueryParam("wikiId"))
	if err != nil {
		log.Printf("failed to convert wikiId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	var message model.MessageToTraQ_POST
	err = c.Bind(&message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	var messageContents model.SodanContent_fromDB
	err = h.db.Select(&messageContents, "select * from messages where wiki_id = ? order by created_at desc limit 1", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get sodanContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	messageID := messageContents.MessageID
	err = h.scraper.MessageToTraQ(message.Content, messageID)

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, h.GetSodanHandler(c))
}

func (h *Handler) PostRepliesToTraQ(c echo.Context) error {
	wikiId, err := strconv.Atoi(c.QueryParam("wikiId"))
	if err != nil {
		log.Printf("failed to convert wikiId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	var message model.MessageToTraQ_POST
	err = c.Bind(&message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	var messageContents model.SodanContent_fromDB
	err = h.db.Select(&messageContents, "select * from messages where wiki_id = ? order by created_at desc limit 1", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get sodanContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	channelID := messageContents.ChannelID
	err = h.scraper.MessageToTraQ(message.Content, channelID)

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, h.GetSodanHandler(c))
}
