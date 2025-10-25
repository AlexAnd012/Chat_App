package data

import "time"

type Message struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	Text     string    `json:"text"`
	Sendtime time.Time `json:"sendtime"`
	IsSystem bool      `json:"isSystem,omitempty"`
}
