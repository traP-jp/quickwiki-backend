package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"strconv"

	"github.com/labstack/echo"
	"github.com/traPtitech/go-traq"
)

func (h *Handler) PostMessageToTraQ(c echo.Context) error {
	var message model.MessageToTraQ_POST
	err := c.Bind(&message)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	owner, err := h.GetUserInfo(c)
	if err != nil {
		log.Println("failed to get user info : ", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get user traQID")
	}

	var resp *traq.Message
	channelId := "aff37b5f-0911-4255-81c3-b49985c8943f" //random/sodan
	resp, err = h.scraper.MessageToTraQ(message.Content, channelId)
	if err != nil {
		log.Printf("MessageToTraQ err : %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "For some reason we could not post the message to traQ")
	}

	messageId := resp.Id
	_, err = h.db.Exec("INSERT INTO anonSodans (message_traq_id,user_traq_id) VALUES (?,?)", messageId, owner.TraqID)
	if err != nil {
		log.Printf("DB Error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	return c.JSON(http.StatusOK, "ok")
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
	err = h.db.Select(&messageContents, "select * from anonSodans where wiki_id = ? order by created_at asc limit 1", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get sodanContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	messageID := messageContents.MessageID
	var anonSodan model.AnonSodans_fromDB
	err = h.db.Select(&anonSodan, "select * from anonSodans where wiki_id = ? order by created_at desc limit 1", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get anonSodanContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	owner, err := h.GetUserInfo(c)
	if err != nil {
		log.Println("failed to get user info : ", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get user traQID")
	}
	if owner.TraqID != anonSodan.UserTraqID {
		return echo.NewHTTPError(http.StatusForbidden, "This is not your MESSAGE!") //これでいいのか?
	}

	err = h.scraper.MessageEditOnTraQ(message.Content, messageID)
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

	resp, err := h.scraper.MessageToTraQ(message.Content, channelID)

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	messageId := resp.Id
	owner, err := h.GetUserInfo(c)
	if err != nil {
		log.Println("failed to get user info : ", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get user traQID")
	}

	_, err = h.db.Exec("INSERT INTO anonSodans (message_traq_id,user_traq_id) VALUES (?,?)", messageId, owner.TraqID)
	if err != nil {
		log.Printf("DB Error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, h.GetSodanHandler(c))
}
