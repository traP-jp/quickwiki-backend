package handler

type LectureFromDB struct {
	ID         int    `db:"id"`
	Title      string `db:"title"`
	Content    string `db:"content"`
	FolderID   int    `db:"folder_id"`
	FolderPath string `db:"folderpath"`
}

type Lecture struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	FolderPath string `json:"folderpath"`
}

type FolderFromDB struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type File struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	IsFolder bool   `json:"isFolder"`
}

type LectureOnlyName struct {
	ID    int    `db:"id"`
	Title string `db:"title"`
}
