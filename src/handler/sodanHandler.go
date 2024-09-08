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

func (h *Handler) PostAnonSodanToTraQ(c echo.Context) error {
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

	if message.isDraft == nil {
		message.isDraft = false
	}

	if message.IsDraft {
		result, err := h.db.Exec("INSERT INTO anonSodanDrafts (content,user_traq_id) VALUES (?,?)", message.Content, owner.Name)
		if err != nil {
			log.Println("failed to insert anonSodanDrafts : ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert anonSodanDrafts")
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Println("failed to get anonSodanDrafts id : ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get anonSodanDrafts id")
		}
		var resp model.AnonSodanDraftRespons
		resp.ID = id
		resp.UserTraqID = owner.TraqID
		resp.Content = message.Content
		return c.JSON(http.StatusOK, resp)
	} else {
		var postResp *traq.Message
		channelId := "aff37b5f-0911-4255-81c3-b49985c8943f" //random/sodan
		postResp, err = h.scraper.MessageToTraQ(message.Content, channelId)
		if err != nil {
			log.Printf("MessageToTraQ err : %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "For some reason we could not post the message to traQ")
		}

		messageId := postResp.Id
		_, err = h.db.Exec("INSERT INTO anonSodans (message_traq_id,user_traq_id) VALUES (?,?)", messageId, owner.TraqID)
		if err != nil {
			log.Printf("DB Error: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}

		if err != nil {
			log.Println("post message err : ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
		}
		var resp model.AnonSodanRespons
		resp.MessageTraqId = messageId
		resp.UserTraqID = owner.TraqID
		resp.Content = message.Content
		return c.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) PatchAnonSodanToTraQ(c echo.Context) error {
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

	owner, err := h.GetUserInfo(c)
	if err != nil {
		log.Println("failed to get user info : ", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get user traQID")
	}

	if message.isDraft == nil {
		message.isDraft = false
	}

	if message.isDraft {
		_, err := h.db.Exec("UPDATE anonSodanDrafts SET content = ? , user_traq_id = ? WHERE id = ?", message.Content, owner.Name, wikiId)
		if err != nil {
			log.Println("failed to insert anonSodanDrafts : ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert anonSodanDrafts")
		}
		var resp model.AnonSodanDraftRespons
		resp.ID = wikiId
		resp.UserTraqID = owner.TraqID
		resp.Content = message.Content
		return c.JSON(http.StatusOK, resp)
	} else {
		var anonSodan model.AnonSodans_fromDB
		err = h.db.GET(&anonSodan, "select * from anonSodans where wiki_id = ? order by created_at asc limit 1", wikiId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("failed to get anonSodanContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		messageID := anonSodan.messageTraqId

		if owner.TraqID != anonSodan.UserTraqID {
			return echo.NewHTTPError(http.StatusForbidden, "This is not your MESSAGE!")
		}

		err = h.scraper.MessageEditOnTraQ(message.Content, messageID)
		if err != nil {
			log.Println("post message err : ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
		}

		return c.JSON(http.StatusOK, h.GetSodanHandler(c))
	}
}

func (h *Handler) PostAnonSodanRepliesToTraQ(c echo.Context) error {
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

	postResp, err := h.scraper.MessageToTraQ(message.Content, channelID)

	if err != nil {
		log.Println("post message err : ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error.")
	}

	messageId := postResp.Id
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

func (h *Handler) GetAnonSodanDraftHandler(c echo.Context) error {
	id := c.QueryParam("messageDraftId")
	owner, err := h.GetUserInfo(c)
	if err != nil {
		log.Println("failed to get user info : ", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "failed to get user traQID")
	}
	if id == "" {
		var anonSodanDrafts []model.AnonSodanDraft_fromDB
		err := h.db.Select(&anonSodanDrafts, "select * from anonSodanDrafts where user_traq_id = ? ", owner.Name)
		if err != nil {
			log.Println("failed to get anonSodanDrafts : err ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get anonSodan drafts")
		}
		var resp []model.AnonSodanDraftRespons
		for _, anonSodanDraft := range anonSodanDrafts {
			var tmp model.anoanonSodanDraftRespons
			tmp.ID = anonSodanDraft.ID
			tmp.Content = anonSodanDraft.Content
			tmp.UserTraqID = anonSodanDraft.UserTraqID
			resp = append(resp, tmp)
		}
		return echo.NewHTTPError(http.StatusOK, resp)
	} else {
		id, err = strconv.Atoi(id)
		if err != nil {
			log.Printf("failed to convert wikiId to int: %v", err)
			return c.JSON(http.StatusBadRequest, err)
		}
		var anonSodanDraftById model.AnonSodanDraft_fromDB
		err = h.db.Select(&anonSodanDraftById, "select * from anonSodanDrafts where id = ? ", id)
		if err != nil {
			log.Println("failed to get anonSodanDrafts : err ", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to get anonSodan drafts")
		}
		var resp []model.AnonSodanDraftRespons
		resp.ID = anonSodanDraftById.ID
		resp.Content = anonSodanDraftById.Content
		resp.UserTraqID = anonSodanDraftById.UserTraqID
		return echo.NewHTTPError(http.StatusOK, resp)
	}
}
