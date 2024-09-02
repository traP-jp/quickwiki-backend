package handler

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"quickwiki-backend/search"
	"strconv"
)

// /sodan/?wikiId=
func (h *Handler) GetSodanHandler(c echo.Context) error {

	Response := model.NewSodanResponse()

	wikiId, err := strconv.Atoi(c.QueryParam("wikiId"))
	if err != nil {
		log.Printf("failed to convert wikiId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	Response.WikiID = wikiId

	var wikiContent model.WikiContent_fromDB
	err = h.db.Get(&wikiContent, "select * from wikis where id = ?", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get wikiContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	Response.Title = wikiContent.Name

	// get tags
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

	// get messages
	var messageContents []model.SodanContent_fromDB
	var howManyMessages int
	err = h.db.Select(&messageContents, "select * from messages where wiki_id = ? order by created_at", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get sodanContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	howManyMessages = len(messageContents)
	log.Println("howManyMessages : ", howManyMessages)
	Response.QuestionMessage.UserTraqID = messageContents[0].UserTraqID
	Response.QuestionMessage.Content = messageContents[0].MessageContent
	Response.QuestionMessage.CreatedAt = messageContents[0].CreatedAt
	Response.QuestionMessage.UpdatedAt = messageContents[0].UpdatedAt
	citedMessagesFromDB := []model.CitedMessage_fromDB{}
	// get citedMessages for question
	err = h.db.Select(&citedMessagesFromDB, "select * from citedMessages where parent_message_id = ?", messageContents[0].ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get citedMessagesFromDB: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	citedMessages := []model.MessageContentForCitations_SodanResponse{}
	for _, citedMessage := range citedMessagesFromDB {
		citedMessageContent := model.MessageContentForCitations_SodanResponse{
			UserTraqID:     citedMessage.UserTraqID,
			CreatedAt:      citedMessage.CreatedAt,
			UpdatedAt:      citedMessage.UpdatedAt,
			MessageContent: citedMessage.Content,
		}
		citedMessages = append(citedMessages, citedMessageContent)
	}
	Response.QuestionMessage.Citations = citedMessages
	for i := 1; i < howManyMessages; i++ {
		ans_Response := model.NewMessageContent_SodanResponse()
		ans_Response.UserTraqID = messageContents[i].UserTraqID
		ans_Response.Content = messageContents[i].MessageContent
		ans_Response.CreatedAt = messageContents[i].CreatedAt
		ans_Response.UpdatedAt = messageContents[i].UpdatedAt
		// get citedMessages for answer
		err = h.db.Select(&citedMessagesFromDB, "select * from citedMessages where parent_message_id = ?", messageContents[i].ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get citedMessagesFromDB: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		citedMessages := []model.MessageContentForCitations_SodanResponse{}
		for _, citedMessage := range citedMessagesFromDB {
			citedMessageContent := model.MessageContentForCitations_SodanResponse{
				UserTraqID:     citedMessage.UserTraqID,
				CreatedAt:      citedMessage.CreatedAt,
				UpdatedAt:      citedMessage.UpdatedAt,
				MessageContent: citedMessage.Content,
			}
			citedMessages = append(citedMessages, citedMessageContent)
		}
		ans_Response.Citations = citedMessages
		Response.AnswerMessages = append(Response.AnswerMessages, *ans_Response)
	}

	// get stamps
	for i := 0; i < howManyMessages; i++ {
		var stamps []model.Stamp_fromDB
		err = h.db.Select(&stamps, "select * from messageStamps where message_id = ?", messageContents[i].ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("failed to get messageStamps: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		howManyStamps := len(stamps)
		log.Println("howManyStamps : ", howManyStamps)
		for j := 0; j < howManyStamps; j++ {
			if i == 0 {
				var stamps_Response model.Stamp_MessageContent
				stamps_Response.StampTraqID = stamps[j].StampTraqID
				stamps_Response.StampCount = stamps[j].StampCount
				Response.QuestionMessage.Stamps = append(Response.QuestionMessage.Stamps, stamps_Response)
			} else {
				var stamps_Response model.Stamp_MessageContent
				stamps_Response.StampTraqID = stamps[j].StampTraqID
				stamps_Response.StampCount = stamps[j].StampCount
				Response.AnswerMessages[i-1].Stamps = append(Response.AnswerMessages[i-1].Stamps, stamps_Response)
			}
		}
	}
	return c.JSON(http.StatusOK, Response)
}

// wikiIdからsearchのResponseを完成させる関数
func WikiIdToResponse(h *Handler, c echo.Context, wikiIds []int) error {
	var Response []model.WikiContentResponse
	for i := 0; i < len(wikiIds); i++ {
		wikiId := wikiIds[i]
		var wikiContent model.WikiContent_fromDB
		tmpSearchContent := model.NewWikiContentResponse()
		err := h.db.Get(&wikiContent, "select * from wikis where id = ?", wikiId)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("failed to get wikiContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		tmpSearchContent.ID = wikiId
		tmpSearchContent.Type = wikiContent.Type
		tmpSearchContent.Title = wikiContent.Name
		tmpSearchContent.Abstract = firstTenChars(wikiContent.Content, 20) //Abstractを入れるべき
		tmpSearchContent.CreatedAt = wikiContent.CreatedAt
		tmpSearchContent.UpdatedAt = wikiContent.UpdatedAt
		tmpSearchContent.OwnerTraqID = wikiContent.OwnerTraqID

		var tags []model.Tag_fromDB
		var howManyTags int
		err = h.db.Select(&tags, "select * from tags where wiki_id = ?", wikiId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		howManyTags = len(tags)
		for j := 0; j < howManyTags; j++ {
			tmpSearchContent.Tags = append(tmpSearchContent.Tags, tags[j].TagName)
		}
		Response = append(Response, *tmpSearchContent)
	}
	return c.JSON(http.StatusOK, Response)
}

// POST/wiki/search の検索はんどら
func (h *Handler) SearchHandler(c echo.Context) error {
	var request model.WikiSearchBody
	err := c.Bind(&request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	// query検索
	var searchResults_Query []int
	if request.Query != "" {
		searchResults_Query = search.Search(request.Query, request.ResultCount, request.From)
		if len(searchResults_Query) == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "No results were found matching that query.")
		}
	}
	// tag検索
	var searchResults_Tags []int
	if len(request.Tags) != 0 {
		log.Println("request.Tags : ", len(request.Tags))
		var searchResultWikiIds [][]int
		for i := 0; i < len(request.Tags); i++ {
			searchResultWikiIds = append(searchResultWikiIds, []int{})
			var tags []model.Tag_fromDB
			err = h.db.Select(&tags, "select * from tags where name = ?", request.Tags[i])
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Printf("failed to get tags: %s\n", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			for _, tag := range tags {
				searchResultWikiIds[i] = append(searchResultWikiIds[i], tag.WikiID)
			}

			var tagLike []model.Tag_fromDB
			err = h.db.Select(&tagLike, "select * from tags where name like ? and name != ?", "%"+request.Tags[i]+"%", request.Tags[i])
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				log.Printf("failed to get tags: %s\n", err)
				return c.NoContent(http.StatusInternalServerError)
			}
			for _, tag := range tagLike {
				searchResultWikiIds[i] = append(searchResultWikiIds[i], tag.WikiID)
			}
		}
		log.Println("searchResultWikiIds : ", searchResultWikiIds)

		intersection := searchResultWikiIds[0]
		unionSet := searchResultWikiIds[0]
		log.Println("len(searchResultWikiIds)", len(searchResultWikiIds))
		for i := 0; i < len(searchResultWikiIds); i++ {
			intersection = intersectUsingMap(intersection, searchResultWikiIds[i])
			unionSet = unionUsingMap(unionSet, searchResultWikiIds[i])
		}
		if len(intersection) == 0 {
			return c.NoContent(http.StatusNotFound)
		}
		if len(unionSet) == 0 {
			return c.NoContent(http.StatusNotFound)
		}
		for _, tmp := range intersection {
			searchResults_Tags = append(searchResults_Tags, tmp) //ここでtagの検索結果の積集合を選択している
		}
		log.Println("searchResults_Tags", searchResults_Tags)
	}

	//検索結果の調整
	var Response_WikiId []int
	if searchResults_Query != nil && searchResults_Tags != nil {
		for _, tmp := range intersectUsingMap(searchResults_Query, searchResults_Tags) {
			Response_WikiId = append(Response_WikiId, tmp) //queryとtagの積集合で検索
		}
	} else {
		for _, tmp := range unionUsingMap(searchResults_Query, searchResults_Tags) {
			Response_WikiId = append(Response_WikiId, tmp) //queryとtagのどちらかしかなかったらそのどちらかで検索
		}
	}

	return WikiIdToResponse(h, c, Response_WikiId)
}

// /wiki/tag?tag=&tag=...のtagからwikiを返すはんどら
func (h *Handler) GetWikiByTagHandler(c echo.Context) error {
	params := c.QueryParams()
	requestTags := params["tag"]
	var searchResults_Tags []int
	if len(requestTags) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request query")
	}
	log.Println("tags : ", len(requestTags))
	var searchResultWikiIds [][]int
	for i, requestTag := range requestTags {
		searchResultWikiIds = append(searchResultWikiIds, []int{})
		var tags []model.Tag_fromDB
		err := h.db.Select(&tags, "select * from tags where name = ?", requestTag)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		for _, tag := range tags {
			searchResultWikiIds[i] = append(searchResultWikiIds[i], tag.WikiID)
		}

		var tagLike []model.Tag_fromDB
		err = h.db.Select(&tagLike, "select * from tags where name like ? and name != ?", "%"+requestTag+"%", requestTag)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Printf("failed to get tags: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		for _, tag := range tagLike {
			searchResultWikiIds[i] = append(searchResultWikiIds[i], tag.WikiID)
		}

	}
	log.Println("searchResultWikiIds : ", searchResultWikiIds)

	intersection := searchResultWikiIds[0]
	unionSet := searchResultWikiIds[0]
	log.Println("len(searchResultWikiIds)", len(searchResultWikiIds))
	for i := 0; i < len(searchResultWikiIds); i++ {
		intersection = intersectUsingMap(intersection, searchResultWikiIds[i])
		unionSet = unionUsingMap(unionSet, searchResultWikiIds[i])
	}
	if len(intersection) == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	if len(unionSet) == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	for _, tmp := range intersection {
		searchResults_Tags = append(searchResults_Tags, tmp) //ここでtagの検索結果の積集合を選択している
	}
	log.Println("searchResults_Tags", searchResults_Tags)

	return WikiIdToResponse(h, c, searchResults_Tags)
}

// /wiki/tag (POST)
func (h *Handler) PostTagHandler(c echo.Context) error {
	tagRequest := model.Tag_Post{}
	err := c.Bind(&tagRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	_, err = h.db.Exec("INSERT INTO tags (name,tag_score,wiki_id) VALUES (?,?,?)", tagRequest.Tag, 1, tagRequest.WikiID)
	if err != nil {
		log.Printf("failed to insert tag: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, tagRequest)
}

// /wiki/tag (PATCH)
func (h *Handler) EditTagHandler(c echo.Context) error {
	tagRequest := model.Tag_Patch{}
	err := c.Bind(&tagRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	_, err = h.db.Exec("UPDATE tags SET name = ? WHERE name = ? AND wiki_id = ?", tagRequest.NewTag, tagRequest.Tag, tagRequest.WikiID)
	if err != nil {
		log.Printf("failed to update tag: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, model.Tag_Post{
		WikiID: tagRequest.WikiID,
		Tag:    tagRequest.NewTag,
	})
}

// /wiki/tag (DELETE)
func (h *Handler) DeleteTagHandler(c echo.Context) error {
	tagRequest := model.Tag_Post{}
	err := c.Bind(&tagRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	_, err = h.db.Exec("DELETE FROM tags WHERE name = ? AND wiki_id = ?", tagRequest.Tag, tagRequest.WikiID)
	if err != nil {
		log.Printf("failed to delete tag: %+v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	return c.JSON(http.StatusOK, tagRequest)
}

// /tag
func (h *Handler) GetTagsHandler(c echo.Context) error {
	tags := []string{}
	err := h.db.Select(&tags, "SELECT DISTINCT name FROM tags")
	if err != nil {
		log.Printf("failed to get tags: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, tags)
}
