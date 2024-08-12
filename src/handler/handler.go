package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) PingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

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

	folderTree := []string{}
	folderName := ""
	err = h.db.Get(&folderName, "SELECT name FROM folders WHERE id = ?", folderID)
	if err != nil {
		log.Printf("failed to get folder name: %v", err)
		return c.JSON(http.StatusInternalServerError, err)
	}
	folderTree = append(folderTree, folderName)
	for {
		parentID := 0
		err = h.db.Get(&parentID, "SELECT parent_id FROM folders WHERE id = ?", folderID)
		if err != nil {
			log.Printf("failed to get parent_id: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}

		// root folder
		if parentID == 0 {
			break
		}

		err = h.db.Get(&folderName, "SELECT name FROM folders WHERE id = ?", parentID)
		if err != nil {
			log.Printf("failed to get folder name: %v", err)
			return c.JSON(http.StatusInternalServerError, err)
		}
		folderTree = append(folderTree, folderName)
		folderID = parentID
	}
	folderPath := ""
	for i, name := range folderTree {
		if i == 0 {
			folderPath = name
		} else {
			folderPath = name + "/" + folderPath
		}
	}
	folderPath = "/" + folderPath

	lecturesWithFolderPath := []Lecture{}
	for _, lecture := range lectures {
		lecturesWithFolderPath = append(lecturesWithFolderPath, Lecture{
			ID:        lecture.ID,
			Title:     lecture.Title,
			Content:   lecture.Content,
			FolderPath: folderPath,
		})
	}

	return c.JSON(http.StatusOK, lecturesWithFolderPath)
}
