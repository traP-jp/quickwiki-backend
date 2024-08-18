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

// /sodan
func (h *Handler) GetSodanHandler(c echo.Context) error {
	wikiId := c.Param("wikiId")
	if wikiId == "" {
		return c.JSON(http.StatusBadRequest, "wikiId is required")
	}

	wiki := Wiki{}
	err := h.db.Get(&wiki, "SELECT * FROM wikis WHERE id = ?", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get wiki: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	questionMessage := Message{}
	err = h.db.Get(&questionMessage, "SELECT * FROM messages WHERE wiki_id = ? ORDER BY created_at DESC LIMIT 1", wikiId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get question message: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	answerMessages := []Message{}
	err = h.db.Select(&answerMessages, "SELECT * FROM messages WHERE wiki_id = ? AND id != ? ORDER BY created_at", wikiId, questionMessage.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get answer messages: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	questionMessageStamp := []Stamp{}
	err = h.db.Select(&questionMessageStamp, "SELECT * FROM messageStamps WHERE message_id = ?", questionMessage.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.JSON(http.StatusNotFound, err)
		}
		log.Printf("failed to get question message stamps: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	questionMessageStampResponse := []StampResponse{}
	for _, stamp := range questionMessageStamp {
		questionMessageStampResponse = append(questionMessageStampResponse, StampResponse{
			StampID: stamp.StampTraqID,
			Count:   stamp.Count,
		})
	}

	questionMessageResponse := MessageResponse{
		UserTraqID: questionMessage.UserTraqID,
		Content:    questionMessage.Content,
		CreatedAt:  questionMessage.CreatedAt,
		UpdatedAt:  questionMessage.UpdatedAt,
		Stamps:     questionMessageStampResponse,
	}

	answerMessagesResponse := []MessageResponse{}
	for _, message := range answerMessages {
		stamps := []Stamp{}
		err = h.db.Select(&stamps, "SELECT * FROM messageStamps WHERE message_id = ?", message.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.JSON(http.StatusNotFound, err)
			}
			log.Printf("failed to get stamps: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		stampsResponse := []StampResponse{}
		for _, stamp := range stamps {
			stampsResponse = append(stampsResponse, StampResponse{
				StampID: stamp.StampTraqID,
				Count:   stamp.Count,
			})
		}

		answerMessagesResponse = append(answerMessagesResponse, MessageResponse{
			UserTraqID: message.UserTraqID,
			Content:    message.Content,
			CreatedAt:  message.CreatedAt,
			UpdatedAt:  message.UpdatedAt,
			Stamps:     stampsResponse,
		})
	}

	res := SodanResponse{
		ID:              wiki.ID,
		Title:           wiki.Name,
		Tags:            []Tag{{Name: "test"}, {Name: "test2"}},
		QuestionMessage: questionMessageResponse,
		AnswerMessages:  answerMessagesResponse,
	}

	return c.JSON(http.StatusOK, res)
}
