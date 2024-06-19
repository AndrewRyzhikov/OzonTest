package entity

import (
	"time"
)

type Post struct {
	PostID    int        `json:"postID" db:"post_id"`
	UserID    int        `json:"userID" db:"user_id"`
	Content   string     `json:"content" db:"content"`
	Comments  []*Comment `json:"comments"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	IsOpen    bool       `json:"isOpen" db:"is_open"`
}

type PostFilter struct {
	UserID int `json:"userID"`
}
