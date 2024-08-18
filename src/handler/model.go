package handler

import "time"

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

type Wiki struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	OwnerTraqID string    `db:"owner_traq_id"`
}

type Message struct {
	ID         int       `db:"id"`
	WikiID     int       `db:"wiki_id"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	UserTraqID string    `db:"user_traq_id"`
	ChannelID  string    `db:"channel_id"`
	MessageID  string    `db:"message_id"`
}

type Stamp struct {
	ID          int    `db:"id"`
	MessageID   int    `db:"message_id"`
	StampTraqID string `db:"stamp_traq_id"`
	Count       int    `db:"count"`
}

type SodanResponse struct {
	ID              int               `json:"id"`
	Title           string            `json:"title"`
	Tags            []Tag             `json:"tags"`
	QuestionMessage MessageResponse   `json:"questionMessage"`
	AnswerMessages  []MessageResponse `json:"answerMessages"`
}

type MessageResponse struct {
	UserTraqID string          `json:"userTraqId"`
	Content    string          `json:"content"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
	Stamps     []StampResponse `json:"stamps"`
}

type StampResponse struct {
	StampID string `json:"stampId"`
	Count   int    `json:"count"`
}

type Tag struct {
	Name string `db:"name"`
}
