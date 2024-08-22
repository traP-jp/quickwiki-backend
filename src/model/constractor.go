package model

import "time"

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

// memoのbodyから情報を取ってくる型のコンストラクタ
func NewGetMemoBody() *GetMemoBody {
	return &GetMemoBody{
		ID:          0,
		Title:       "",
		Content:     "",
	}
}