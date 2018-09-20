package app

import "time"

type Note struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`

	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
type Notes []Note
