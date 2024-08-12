package handler

type LectureFromDB struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	Content   string `db:"content"`
	FolderID  int    `db:"folder_id"`
}

type Lecture struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	FolderPath string `json:"folderpath"`
}