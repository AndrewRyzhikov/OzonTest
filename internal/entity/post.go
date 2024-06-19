package entity

import (
	"time"
)

type Post struct {
	PostID    int        `json:"postID"`
	UserID    int        `json:"userID"`
	Content   string     `json:"content"`
	Comments  []*Comment `json:"comments"`
	Timestamp time.Time  `json:"timestamp"`
	IsOpen    bool       `json:"isOpen"`
}

type PostFilter struct {
	UserID int `json:"userID"`
}
