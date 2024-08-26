package handler

import (
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"quickwiki-backend/scraper"
	"strconv"
	"strings"
	"time"

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

// /lecture/byFolder/id/:folderId
func (h *Handler) GetLectureByFolderIDHandler(c echo.Context) error {
	folderID, err := strconv.Atoi(c.Param("folderId"))
	if err != nil {
		log.Printf("failed to convert folderId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	lectures := []model.LectureFromDB{}
	err = h.db.Select(&lectures, "SELECT * FROM lectures WHERE folder_id = ?", folderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get lectures: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	lecturesWithFolderPath := []model.Lecture{}
	for _, lecture := range lectures {
		lecturesWithFolderPath = append(lecturesWithFolderPath, model.Lecture{
			ID:         lecture.ID,
			Title:      lecture.Title,
			Content:    lecture.Content,
			FolderPath: lecture.FolderPath,
		})
	}

	return c.JSON(http.StatusOK, lecturesWithFolderPath)
}

// /lecture/byFolder/path
func (h *Handler) GetLectureByFolderPathHandler(c echo.Context) error {
	folderPath := c.QueryParam("folderPath")

	folderPath = "/" + strings.ReplaceAll(folderPath, "-", " /")
	lectures := []model.LectureFromDB{}
	err := h.db.Select(&lectures, "SELECT * FROM lectures WHERE folder_path = ?", folderPath)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get lectures: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	lecturesWithFolderPath := []model.Lecture{}
	for _, lecture := range lectures {
		lecturesWithFolderPath = append(lecturesWithFolderPath, model.Lecture{
			ID:         lecture.ID,
			Title:      lecture.Title,
			Content:    lecture.Content,
			FolderPath: lecture.FolderPath,
		})
	}

	return c.JSON(http.StatusOK, lecturesWithFolderPath)
}

// /lecture/folder/:folderId
func (h *Handler) GetLectureChildFolderHandler(c echo.Context) error {
	folderID, err := strconv.Atoi(c.Param("folderId"))
	if err != nil {
		log.Printf("failed to convert folderId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	files := []model.File{}

	childFolders := []model.FolderFromDB{}
	err = h.db.Select(&childFolders, "SELECT * FROM folders WHERE parent_id = ?", folderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get child folders: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, folder := range childFolders {
		files = append(files, model.File{
			ID:       folder.ID,
			Name:     folder.Name,
			IsFolder: true,
		})
	}

	childLectures := []model.LectureOnlyName{}
	err = h.db.Select(&childLectures, "SELECT id, title FROM lectures WHERE folder_id = ?", folderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get child lectures: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, lecture := range childLectures {
		files = append(files, model.File{
			ID:       lecture.ID,
			Name:     lecture.Title,
			IsFolder: false,
		})
	}

	return c.JSON(http.StatusOK, files)
}

// /lecture/:lectureId
func (h *Handler) GetLectureHandler(c echo.Context) error {
	lectureID, err := strconv.Atoi(c.Param("lectureId"))
	if err != nil {
		log.Printf("failed to convert lectureId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	lecture := model.LectureFromDB{}
	err = h.db.Get(&lecture, "SELECT * FROM lectures WHERE id = ?", lectureID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get lecture: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, model.Lecture{
		ID:         lecture.ID,
		Title:      lecture.Title,
		Content:    lecture.Content,
		FolderPath: lecture.FolderPath,
	})
}

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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
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
				Response.AnswerMessages[i-1].Stamps = append(Response.QuestionMessage.Stamps, stamps_Response)
			}
		}
	}
	return c.JSON(http.StatusOK, Response)
}

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
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get wikiContent: %s\n", err)
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get tags: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	howManyTags = len(tags)
	for i := 0; i < howManyTags; i++ {
		Response.Tags = append(Response.Tags, tags[i].TagName)
	}

	return c.JSON(http.StatusOK, Response)
}

func (h *Handler) GetFileHandler(c echo.Context) error {
	fileID := c.Param("fileId")

	resp, err := h.scraper.GetFile(fileID)
	if err != nil {
		log.Printf("failed to get file: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := c.Response()
	response.Header().Set("Cache-Control", "no-cache")
	response.Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	response.Header().Set(echo.HeaderAccessControlExposeHeaders, "Content-Disposition")
	response.Header().Set(echo.HeaderContentDisposition, "attachment; filename="+fileID)
	response.WriteHeader(http.StatusOK)
	io.Copy(response.Writer, resp.Body)
	return c.NoContent(http.StatusOK)
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
	copy(Response.Tags, getMemoBody.Tags)

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
		err = h.db.Get(&wikiContent, "select owner_traq_id from wikis where id = ?", getMemoBody.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("failed to get wikiContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if wikiContent.OwnerTraqID != owner.TraqID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You are not the owner of this memo.")
		}
		now := time.Now()
		_, err := h.db.Exec("UPDATE wikis SET name = ?,updated_at = ?,content = ? where wiki_id = ?", getMemoBody.Title, now, getMemoBody.Content, getMemoBody.ID)
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
		err = h.db.Get(&tags, "select * from tags where wiki_id = ?", getMemoBody.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("failed to get wikiContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		var resTags []string
		for i := 0; i < len(tags); i++ {
			resTags[i] = tags[i].TagName
		}
	}

	return c.JSON(http.StatusOK, Response)
}

// DELETE/memo のハンドラー
func (h *Handler) DeleteMemoHandler(c echo.Context) error {

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
			log.Printf("failed to get wikiContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		if wikiContent.OwnerTraqID != owner.TraqID {
			return echo.NewHTTPError(http.StatusUnauthorized, "You are not the owner of this memo.")
		}
		_, err := h.db.Exec("DELETE from wikis where wiki_id = ?", getMemoBody.ID)
		if err != nil {
			log.Printf("DB Error: %s", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
		Response.WikiID = wikiContent.ID
		Response.Title = wikiContent.Name
		Response.Content = wikiContent.Content
		Response.OwnerTraqID = wikiContent.OwnerTraqID
		Response.CreatedAt = wikiContent.CreatedAt
		Response.UpdatedAt = wikiContent.UpdatedAt

		var tags []model.Tag_fromDB
		err = h.db.Get(&tags, "select * from tags where wiki_id = ?", getMemoBody.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			log.Printf("failed to get wikiContent: %s\n", err)
			return c.NoContent(http.StatusInternalServerError)
		}
		var resTags []string
		for i := 0; i < len(tags); i++ {
			resTags[i] = tags[i].TagName
		}
	}

	return c.JSON(http.StatusOK, Response)
}

// /wiki/tag
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

// /me
func (h *Handler) GetMeHandler(c echo.Context) error {
	user, err := h.GetUserInfo(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, user)
}

// /lecture
func (h *Handler) PostLectureHandler(c echo.Context) error {
	lecturePost := model.Lecture_Post{}
	err := c.Bind(&lecturePost)
	if err != nil {
		log.Printf("failed to bind lecturePost: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	folderTree := strings.Split(lecturePost.FolderPath, "/")
	folderTree[0] = "root"

	// insert folder if not exists
	for i, folder := range folderTree {
		if i == 0 {
			continue
		}

		err := h.db.Get(&folder, "SELECT id FROM folders WHERE name = ? AND parent_id = ?", folder, folderTree[i-1])
		if errors.Is(err, sql.ErrNoRows) {
			_, err = h.db.Exec("INSERT INTO folders (name, parent_id) VALUES (?, ?)", folder, folderTree[i-1])
			if err != nil {
				log.Printf("failed to insert folder: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
		} else if err != nil {
			log.Printf("failed to get folder: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		}
	}

	// insert lecture
	folderID := 0
	err = h.db.Get(&folderID, "SELECT id FROM folders WHERE name = ? AND parent_id = ?", folderTree[len(folderTree)-1], folderTree[len(folderTree)-2])
	if err != nil {
		log.Printf("failed to get folder: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	result, err := h.db.Exec("INSERT INTO lectures (title, content, folder_path, folder_id) VALUES (?, ?, ?, ?)",
		lecturePost.Title, lecturePost.Content, lecturePost.FolderPath, folderID)
	if err != nil {
		log.Printf("failed to insert lecture: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
	}

	lectureID, _ := result.LastInsertId()
	res := model.Lecture{
		ID:         int(lectureID),
		Title:      lecturePost.Title,
		Content:    lecturePost.Content,
		FolderPath: lecturePost.FolderPath,
	}

	return c.JSON(http.StatusOK, res)
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
