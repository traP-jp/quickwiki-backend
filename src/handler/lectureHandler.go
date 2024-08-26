package handler

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo"
	"log"
	"net/http"
	"quickwiki-backend/model"
	"strconv"
	"strings"
)

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
	folderTreeIds := []int{1}

	// insert folder if not exists
	for i, folder := range folderTree {
		if i == 0 {
			continue
		}

		id := 0
		err := h.db.Get(&id, "SELECT id FROM folders WHERE name = ? AND parent_id = ?", folder, folderTreeIds[i-1])
		if errors.Is(err, sql.ErrNoRows) {
			res, err := h.db.Exec("INSERT INTO folders (name, parent_id) VALUES (?, ?)", folder, folderTreeIds[i-1])
			if err != nil {
				log.Printf("failed to insert folder: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
			folderId, err := res.LastInsertId()
			if err != nil {
				log.Printf("failed to get folderId: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
			folderTreeIds = append(folderTreeIds, int(folderId))
		} else if err != nil {
			log.Printf("failed to get folder: %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
		} else {
			folderId := 0
			err := h.db.Get(&folderId, "SELECT id FROM folders WHERE name = ? AND parent_id = ?", folder, folderTreeIds[i-1])
			if err != nil {
				log.Printf("failed to get folder: %v", err)
				return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
			}
			folderTreeIds = append(folderTreeIds, folderId)
		}
	}

	// insert lecture
	folderID := 0
	err = h.db.Get(&folderID, "SELECT id FROM folders WHERE name = ? AND parent_id = ?", folderTree[len(folderTree)-1], folderTreeIds[len(folderTree)-2])
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
