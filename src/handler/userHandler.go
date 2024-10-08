package handler

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"strconv"
)

// /me
func (h *Handler) GetMeHandler(c echo.Context) error {
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

	var wikis []model.WikiContent_fromDB
	err = h.db.Select(&wikis, "SELECT * FROM wikis WHERE owner_traq_id = ?", user.TraqID)
	if err != nil {
		log.Printf("failed to get wikis: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var res []model.WikiContentResponse
	for _, wiki := range wikis {
		// get tag
		var tags []model.Tag_fromDB
		err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		var tagsRes []string
		for _, tag := range tags {
			tagsRes = append(tagsRes, tag.TagName)
		}

		var favorites int
		err = h.db.Get(&favorites, "SELECT COUNT(*) FROM favorites WHERE wiki_id = ?", wiki.ID)
		if err != nil {
			log.Printf("failed to get favorites: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
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
			Favorites:   favorites,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// /wiki/user/favoite
func (h *Handler) GetUserFavoriteWikiHandler(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	var wikis []model.WikiContent_fromDB
	err = h.db.Select(&wikis, "SELECT * FROM wikis WHERE id IN (SELECT wiki_id FROM favorites WHERE user_traq_id = ?)", user.TraqID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get wikis: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	var res []model.WikiContentResponse
	for _, wiki := range wikis {
		// get tag
		var tags []model.Tag_fromDB
		err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		var tagsRes []string
		for _, tag := range tags {
			tagsRes = append(tagsRes, tag.TagName)
		}

		var favorites int
		err = h.db.Get(&favorites, "SELECT COUNT(*) FROM favorites WHERE wiki_id = ?", wiki.ID)
		if err != nil {
			log.Printf("failed to get favorites: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
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
			Favorites:   favorites,
		})
	}

	return c.JSON(http.StatusOK, res)
}

// /wiki/user/favoite POST
func (h *Handler) PostUserFavoriteWikiHandler(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	wikiIDPost := struct {
		WikiID string `json:"wikiId"`
	}{}
	err = c.Bind(&wikiIDPost)
	if err != nil {
		log.Printf("failed to bind wikiID: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	wikiID, err := strconv.Atoi(wikiIDPost.WikiID)
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
		log.Printf("wikiid: %d", wikiID)
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
	var tags []model.Tag_fromDB
	err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get tags: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	var tagsRes []string
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
		Favorites:   0,
	}

	return c.JSON(http.StatusOK, res)
}

// /wiki/user/favoite DELETE
func (h *Handler) DeleteUserFavoriteWikiHandler(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}

	wikiIDPost := struct {
		WikiID string `json:"wikiId"`
	}{}
	err = c.Bind(&wikiIDPost)
	if err != nil {
		log.Printf("failed to bind wikiID: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	wikiID, err := strconv.Atoi(wikiIDPost.WikiID)
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
	var tags []model.Tag_fromDB
	err = h.db.Select(&tags, "SELECT * FROM tags WHERE wiki_id = ?", wiki.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get tags: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	var tagsRes []string
	for _, tag := range tags {
		tagsRes = append(tagsRes, tag.TagName)
	}

	var favorites int
	err = h.db.Get(&favorites, "SELECT COUNT(*) FROM favorites WHERE wiki_id = ?", wiki.ID)
	if err != nil {
		log.Printf("failed to get favorites: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
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
		Favorites:   favorites,
	}

	_, err = h.db.Exec("DELETE FROM favorites WHERE user_traq_id = ? AND wiki_id = ?", user.TraqID, wikiID)
	if err != nil {
		log.Printf("failed to delete favorite: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, res)
}
