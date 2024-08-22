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

// sqlのwikisから情報を取ってくるときに使う
type WikiContent_fromDB struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	OwnerTraqID string    `db:"owner_traq_id"`
	Content     string    `db:"content"`
}

// sqlよりtagを取ってくるときに使う
type Tag_fromDB struct {
	TagID    int     `db:"id"`
	WikiID   int     `db:"wiki_id"`
	TagName  string  `db:"name"`
	TagScore float64 `db:"tag_score"`
}

// sqlのmessagesから情報を取ってくるときに使う
type SodanContent_fromDB struct {
	ID             int       `db:"id"`
	WikiID         int       `db:"wiki_id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	UserTraqID     string    `db:"user_traq_id"`
	MessageID      string    `db:"message_traq_id"`
	ChannelID      string    `db:"channel_id"`
	MessageContent string    `db:"content"`
}

// sqlのmessageStampsから情報を取ってくるときに使う
type Stamp_fromDB struct {
	ID          int    `db:"id"`
	MessageID   int    `db:"message_id"`
	StampTraqID string `db:"stamp_traq_id"`
	StampCount  int    `db:"count"`
}

// SodanResponseで使うstampの構造体
type Stamp_MessageContent struct {
	StampTraqID string `json:"stampId"`
	StampCount  int    `json:"count"`
}

// SodanResponseで使うMessageの構造体
type MessageContent_SodanResponse struct {
	UserTraqID string                                     `json:"userTraqId"`
	Content    string                                     `json:"content"`
	CreatedAt  time.Time                                  `json:"createdAt"`
	UpdatedAt  time.Time                                  `json:"updatedAt"`
	Stamps     []Stamp_MessageContent                     `json:"stamps"`
	Citations  []MessageContentForCitations_SodanResponse `json:"citations"`
}

type MessageContentForCitations_SodanResponse struct {
	UserTraqID     string    `json:"userTraqId"`
	MessageContent string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// MessageContent_SodanResponseのコンストラクタ関数
func NewMessageContent_SodanResponse() *MessageContent_SodanResponse {
	return &MessageContent_SodanResponse{
		UserTraqID: "",
		Content:    "",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Stamps:     []Stamp_MessageContent{},
	}
}

// /sodan?wikiId= の Response
type SodanResponse struct {
	WikiID          int                            `json:"id"`
	Title           string                         `json:"title"`
	Tags            []string                       `json:"tags"`
	QuestionMessage MessageContent_SodanResponse   `json:"questionMessage"`
	AnswerMessages  []MessageContent_SodanResponse `json:"answerMessages"`
}

// SodanResponseのコンストラクタ関数
func NewSodanResponse() *SodanResponse {
	return &SodanResponse{
		WikiID:          0,
		Title:           "",
		Tags:            []string{},
		QuestionMessage: MessageContent_SodanResponse{},
		AnswerMessages:  []MessageContent_SodanResponse{},
	}
}

// sqlのmemosより情報を取ってくるときに使う
type MemoContent_fromDB struct {
	ID          int       `db:"id"`
	WikiID      int       `db:"wiki_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	OwnerTraqID string    `db:"owner_traq_id"`
	Content     string    `db:"content"`
}

// GET/memo?wikiId のResponse構造体
type MemoResponse struct {
	WikiID      int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	OwnerTraqID string    `json:"ownerTraqId"`
	Tags        []string  `json:"tags"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// MemoResponseのコンストラクタ関数
func NewMemoResponse() *MemoResponse {
	return &MemoResponse{
		WikiID:      0,
		Title:       "",
		Content:     "",
		OwnerTraqID: "",
		Tags:        []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type CitedMessage_fromDB struct {
	ID              int       `db:"id"`
	ParentMessageID int       `db:"parent_message_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	UserTraqID      string    `db:"user_traq_id"`
	MessageTraqID   string    `db:"message_traq_id"`
	ChannelID       string    `db:"channel_id"`
	Content         string    `db:"content"`
}
