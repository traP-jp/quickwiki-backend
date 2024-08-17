package main

import "time"

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
