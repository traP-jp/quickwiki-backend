package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
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

	lectures := []LectureFromDB{}
	err = h.db.Select(&lectures, "SELECT * FROM lectures WHERE folder_id = ?", folderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get lectures: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	lecturesWithFolderPath := []Lecture{}
	for _, lecture := range lectures {
		lecturesWithFolderPath = append(lecturesWithFolderPath, Lecture{
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
	lectures := []LectureFromDB{}
	err := h.db.Select(&lectures, "SELECT * FROM lectures WHERE folderpath = ?", folderPath)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get lectures: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	lecturesWithFolderPath := []Lecture{}
	for _, lecture := range lectures {
		lecturesWithFolderPath = append(lecturesWithFolderPath, Lecture{
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

	files := []File{}

	childFolders := []FolderFromDB{}
	err = h.db.Select(&childFolders, "SELECT * FROM folders WHERE parent_id = ?", folderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get child folders: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, folder := range childFolders {
		files = append(files, File{
			ID:       folder.ID,
			Name:     folder.Name,
			IsFolder: true,
		})
	}

	childLectures := []LectureOnlyName{}
	err = h.db.Select(&childLectures, "SELECT id, title FROM lectures WHERE folder_id = ?", folderID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("failed to get child lectures: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	for _, lecture := range childLectures {
		files = append(files, File{
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

	lecture := LectureFromDB{}
	err = h.db.Get(&lecture, "SELECT * FROM lectures WHERE id = ?", lectureID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get lecture: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, Lecture{
		ID:         lecture.ID,
		Title:      lecture.Title,
		Content:    lecture.Content,
		FolderPath: lecture.FolderPath,
	})
}

// /sodan/?wikiId=
func (h *Handler) GetSodanHandler(c echo.Context) error {

	Response := NewSodanResponse()

	wikiId, err := strconv.Atoi(c.QueryParam("wikiId"))
	if err != nil {
		log.Printf("failed to convert wikiId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}
	Response.WikiID = wikiId

	var wikiContent WikiContent_fromDB
	err = h.db.Get(&wikiContent, "select * from wikis where id = ?", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get wikiContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	Response.Title = wikiContent.Name

	var tags []Tag_fromDB
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

	var messageContents []SodanContent_fromDB
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
	for i := 1; i < howManyMessages; i++ {
		ans_Response := NewMessageContent_SodanResponse()
		ans_Response.UserTraqID = messageContents[i].UserTraqID
		ans_Response.Content = messageContents[i].MessageContent
		ans_Response.CreatedAt = messageContents[i].CreatedAt
		ans_Response.UpdatedAt = messageContents[i].UpdatedAt
		Response.AnswerMessages = append(Response.AnswerMessages, *ans_Response)
	}

	for i := 0; i < howManyMessages; i++ {
		var stamps []Stamp_fromDB
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
				var stamps_Response Stamp_MessageContent
				stamps_Response.StampTraqID = stamps[j].StampTraqID
				stamps_Response.StampCount = stamps[j].StampCount
				Response.QuestionMessage.Stamps = append(Response.QuestionMessage.Stamps, stamps_Response)
			} else {
				var stamps_Response Stamp_MessageContent
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

	Response := NewMemoResponse()

	wikiId, err := strconv.Atoi(c.QueryParam("wikiId"))
	if err != nil {
		log.Printf("failed to convert wikiId to int: %v", err)
		return c.JSON(http.StatusBadRequest, err)
	}

	var memoContent MemoContent_fromDB
	err = h.db.Get(&memoContent, "select * from memos where wiki_id = ?", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNotFound)
		}
		log.Printf("failed to get wikiContent: %s\n", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	Response.WikiID = memoContent.ID
	Response.Title = memoContent.Title
	Response.Content = memoContent.Content
	Response.OwnerTraqID = memoContent.OwnerTraqID
	Response.CreatedAt = memoContent.CreatedAt
	Response.UpdatedAt = memoContent.UpdatedAt

	var tags []Tag_fromDB
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
