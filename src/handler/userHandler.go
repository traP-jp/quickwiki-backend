package handler

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"os"
	"quickwiki-backend/model"
	"strconv"
)

// /me
func (h *Handler) GetMeHandler(c echo.Context) error {
	if os.Getenv("DEV_MODE") == "true" {
		return c.JSON(http.StatusOK, model.Me_Response{
			TraqID:      "kavos",
			DisplayName: "kavos",
			IconUri:     "https://q.trap.jp/api/v3/public/icon/kavos",
		})
	}
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, user)
}

// /wiki/user
func (h *Handler) GetUserWikiHandelr(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	wikis := []model.WikiContent_fromDB{}
	err = h.db.Select(&wikis, "SELECT * FROM wikis WHERE owner_traq_id = ?", user.TraqID)
	if err != nil {
		log.Printf("failed to get wikis: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	res := []model.WikiContentResponse{}
	for _, wiki := range wikis {
		// get tag
		tags := []model.Tag_fromDB{}
		err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		tagsRes := []string{}
		for _, tag := range tags {
			tagsRes = append(tagsRes, tag.TagName)
		}

		res = append(res, model.WikiContentResponse{
			ID:          wiki.ID,
			Type:        wiki.Type,
			Title:       wiki.Name,
			Abstract:    firstTenChars(wiki.Content, 50),
			CreatedAt:   wiki.CreatedAt,
			UpdatedAt:   wiki.UpdatedAt,
			OwnerTraqID: wiki.OwnerTraqID,
			Tags:        tagsRes,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// /wiki/user/favoite
func (h *Handler) GetUserFavoriteWikiHandelr(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	wikis := []model.WikiContent_fromDB{}
	err = h.db.Select(&wikis, "SELECT * FROM wikis WHERE id IN (SELECT wiki_id FROM favorites WHERE user_traq_id = ?)", user.TraqID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get wikis: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	res := []model.WikiContentResponse{}
	for _, wiki := range wikis {
		// get tag
		tags := []model.Tag_fromDB{}
		err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		tagsRes := []string{}
		for _, tag := range tags {
			tagsRes = append(tagsRes, tag.TagName)
		}

		res = append(res, model.WikiContentResponse{
			ID:          wiki.ID,
			Type:        wiki.Type,
			Title:       wiki.Name,
			Abstract:    firstTenChars(wiki.Content, 50),
			CreatedAt:   wiki.CreatedAt,
			UpdatedAt:   wiki.UpdatedAt,
			OwnerTraqID: wiki.OwnerTraqID,
			Tags:        tagsRes,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// /wiki/user/favoite POST
func (h *Handler) PostUserFavoriteWikiHandelr(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	wikiIDStr := ""
	err = c.Bind(&wikiIDStr)
	if err != nil {
		log.Printf("failed to bind wikiID: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	wikiID, err := strconv.Atoi(wikiIDStr)
	if err != nil {
		log.Printf("failed to convert wikiID to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	wikiCount := 0
	err = h.db.Get(&wikiCount, "SELECT COUNT(*) FROM wikis WHERE id = ?", wikiID)
	if err != nil {
		log.Printf("failed to get wiki: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if wikiCount == 0 {
		return c.JSON(http.StatusNotFound, "wiki not found")
	}

	_, err = h.db.Exec("INSERT INTO favorites (user_traq_id, wiki_id) VALUES (?, ?)", user.TraqID, wikiID)
	if err != nil {
		log.Printf("failed to insert favorite: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// set response
	wiki := model.WikiContent_fromDB{}
	err = h.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiID)
	if err != nil {
		log.Printf("failed to get wiki: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// get tag
	tags := []model.Tag_fromDB{}
	err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get tags: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	tagsRes := []string{}
	for _, tag := range tags {
		tagsRes = append(tagsRes, tag.TagName)
	}

	res := model.WikiContentResponse{
		ID:          wiki.ID,
		Type:        wiki.Type,
		Title:       wiki.Name,
		Abstract:    firstTenChars(wiki.Content, 50),
		CreatedAt:   wiki.CreatedAt,
		UpdatedAt:   wiki.UpdatedAt,
		OwnerTraqID: wiki.OwnerTraqID,
		Tags:        tagsRes,
	}

	return c.JSON(http.StatusOK, res)
}

// /wiki/user/favoite DELETE
func (h *Handler) DeleteUserFavoriteWikiHandelr(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	wikiIDStr := ""
	err = c.Bind(&wikiIDStr)
	if err != nil {
		log.Printf("failed to bind wikiID: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	wikiID, err := strconv.Atoi(wikiIDStr)
	if err != nil {
		log.Printf("failed to convert wikiID to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	wikiCount := 0
	err = h.db.Get(&wikiCount, "SELECT COUNT(*) FROM wikis WHERE id = ?", wikiID)
	if err != nil {
		log.Printf("failed to get wiki: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	if wikiCount == 0 {
		return c.JSON(http.StatusNotFound, "wiki not found")
	}

	// set response
	wiki := model.WikiContent_fromDB{}
	err = h.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiID)
	if err != nil {
		log.Printf("failed to get wiki: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	// get tag
	tags := []model.Tag_fromDB{}
	err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get tags: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	tagsRes := []string{}
	for _, tag := range tags {
		tagsRes = append(tagsRes, tag.TagName)
	}

	res := model.WikiContentResponse{
		ID:          wiki.ID,
		Type:        wiki.Type,
		Title:       wiki.Name,
		Abstract:    firstTenChars(wiki.Content, 50),
		CreatedAt:   wiki.CreatedAt,
		UpdatedAt:   wiki.UpdatedAt,
		OwnerTraqID: wiki.OwnerTraqID,
		Tags:        tagsRes,
	}

	_, err = h.db.Exec("DELETE FROM favorites WHERE user_traq_id = ? AND wiki_id = ?", user.TraqID, wikiID)
	if err != nil {
		log.Printf("failed to delete favorite: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}
