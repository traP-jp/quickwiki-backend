package handler

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"strconv"
	"time"
)

// /memo?wikiId=
func (h *Handler) GetMemoHandler(c echo.Context) error {

	Response := model.NewMemoResponse()

	wikiId, err := strconv.Atoi(c.QueryParam("wikiId"))
	if err != nil {
		log.Printf("failed to convert wikiId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	var wikiContent model.WikiContent_fromDB
	err = h.db.Get(&wikiContent, "select * from wikis where id = ?", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("wikiid: %d does not exist.", wikiId)
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("[in get memo]failed to get wikiContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if wikiContent.Type == "memo" {
		Response.WikiID = wikiContent.ID
		Response.Title = wikiContent.Name
		Response.Content = wikiContent.Content
		Response.OwnerTraqID = wikiContent.OwnerTraqID
		Response.CreatedAt = wikiContent.CreatedAt
		Response.UpdatedAt = wikiContent.UpdatedAt
	} else {
		log.Printf("This wikiId exists, but it is not a 'memo'.")
		return c.NoContent(http.StatusNotFound)
	}

	var tags []model.Tag_fromDB
	var howManyTags int
	err = h.db.Select(&tags, "select * from tags where wiki_id = ?", wikiId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get tags: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	howManyTags = len(tags)
	for i := 0; i < howManyTags; i++ {
		Response.Tags = append(Response.Tags, tags[i].TagName)
	}

	return c.JSON(http.StatusOK, Response)
}

// POST/memoのハンドラー
func (h *Handler) PostMemoHandler(c echo.Context) error {

	Response := model.NewMemoResponse()

	getMemoBody := model.NewGetMemoBody()
	err := c.Bind(&getMemoBody)
	if err != nil {
		if getMemoBody.ID != 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
		}
	}
	Response.Title = getMemoBody.Title
	Response.Content = getMemoBody.Content

	owner, err := h.GetUserInfo(c)
	Response.OwnerTraqID = owner.TraqID

	now := time.Now()
	result, err := h.db.Exec("INSERT INTO wikis (name,type,created_at,updated_at,content,owner_traq_id) VALUES (?, ?, ?, ?,?,?)", getMemoBody.Title, "memo", now, now, getMemoBody.Content, owner.TraqID)
	if err != nil {
		log.Printf("DB Error: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	WikiID, _ := result.LastInsertId()
	Response.WikiID = int(WikiID)

	howManyTags := len(getMemoBody.Tags)
	for i := 0; i < howManyTags; i++ {
		_, err = h.db.Exec("INSERT INTO tags (name,tag_score,wiki_id) VALUES (?,?,?)", getMemoBody.Tags[i], 1, WikiID)
		if err != nil {
			log.Printf("DB Error: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}
	for i := 0; i < howManyTags; i++ {
		Response.Tags = append(Response.Tags, getMemoBody.Tags[i])
	}

	Response.CreatedAt = now
	Response.UpdatedAt = now

	return c.JSON(http.StatusOK, Response)
}

// PATCH/memo のはんどらー
func (h *Handler) PatchMemoHandler(c echo.Context) error {

	Response := model.NewMemoResponse()

	getMemoBody := model.NewGetMemoBody()
	err := c.Bind(&getMemoBody)
	if err != nil {
		if getMemoBody.ID != 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
		}
	}

	owner, err := h.GetUserInfo(c)

	if getMemoBody.ID != 0 {
		var wikiContent model.WikiContent_fromDB
		wikiContent.OwnerTraqID = ""
		err = h.db.Get(&wikiContent, "select * from wikis where id = ?", getMemoBody.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("[in patch memo]failed to get wikiContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if wikiContent.OwnerTraqID != owner.TraqID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You are not the owner of this memo.")
		}
		now := time.Now()
		_, err := h.db.Exec("UPDATE wikis SET name = ?,updated_at = ?,content = ? where id = ?", getMemoBody.Title, now, getMemoBody.Content, getMemoBody.ID)
		if err != nil {
			log.Printf("DB Error: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		Response.WikiID = getMemoBody.ID
		Response.Title = getMemoBody.Title
		Response.Content = getMemoBody.Content
		Response.OwnerTraqID = owner.TraqID
		Response.CreatedAt = wikiContent.CreatedAt
		Response.UpdatedAt = now

		var tags []model.Tag_fromDB
		err = h.db.Select(&tags, "select * from tags where wiki_id = ?", getMemoBody.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("[in patch memo]failed to get tag: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		var resTags []string
		for _, tag := range tags {
			resTags = append(resTags, tag.TagName)
		}
	}

	return c.JSON(http.StatusOK, Response)
}

// DELETE/memo のハンドラー
func (h *Handler) DeleteMemoHandler(c echo.Context) error {

	Response := model.NewMemoResponse()

	wikiIDPost := struct {
		ID string `json:"wikiId"`
	}{}
	err := c.Bind(&wikiIDPost)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}
	wikiID, err := strconv.Atoi(wikiIDPost.ID)
	if err != nil {
		log.Printf("failed to convert wikiID to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	owner, err := h.GetUserInfo(c)

	var wikiContent model.WikiContent_fromDB
	wikiContent.OwnerTraqID = ""
	err = h.db.Get(&wikiContent, "select * from wikis where id = ?", wikiID)
	log.Printf("wikiId: %v, wikiContent: %v", wikiID, wikiContent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("[in delete memo]failed to get wikiContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if wikiContent.OwnerTraqID != owner.TraqID {
		return echo.NewHTTPError(http.StatusUnauthorized, "You are not the owner of this memo.")
	}

	// delete tags
	_, err = h.db.Exec("DELETE from tags where wiki_id = ?", wikiID)
	if err != nil {
		log.Printf("failed to delete tags: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	_, err = h.db.Exec("DELETE from wikis where id = ?", wikiID)
	if err != nil {
		log.Printf("failed to delete memo: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}
	Response.WikiID = wikiContent.ID
	Response.Title = wikiContent.Name
	Response.Content = wikiContent.Content
	Response.OwnerTraqID = wikiContent.OwnerTraqID
	Response.CreatedAt = wikiContent.CreatedAt
	Response.UpdatedAt = wikiContent.UpdatedAt

	var tags []model.Tag_fromDB
	err = h.db.Select(&tags, "select * from tags where wiki_id = ?", wikiID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("[in delete memo]failed to get tag: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	var resTags []string
	for i := 0; i < len(tags); i++ {
		resTags[i] = tags[i].TagName
	}

	return c.JSON(http.StatusOK, Response)
}
