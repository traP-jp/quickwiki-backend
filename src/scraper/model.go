package scraper

import "time"

type Wiki_fromDB struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Type        string    `db:"type"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	OwnerTraqID string    `db:"owner_traq_id"`
}

type Message_fromDB struct {
	ID         int       `db:"id"`
	WikiID     int       `db:"wiki_id"`
	Content    string    `db:"content"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	UserTraqID string    `db:"user_traq_id"`
	ChannelID  string    `db:"channel_id"`
	MessageID  string    `db:"message_traq_id"`
}

type Stamp_fromDB struct {
	ID          int    `db:"id"`
	MessageID   int    `db:"message_id"`
	StampTraqID string `db:"stamp_traq_id"`
	Count       int    `db:"count"`
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
