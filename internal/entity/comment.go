package entity

import (
	"time"
)

type Comment struct {
	ID        int        `json:"ID"`
	ParentID  int        `json:"parentID"`
	PostID    int        `json:"postID"`
	UserID    int        `json:"userID"`
	Content   string     `json:"content"`
	Timestamp time.Time  `json:"timestamp"`
	Replies   []*Comment `json:"replies"`
}

type CommentFilter struct {
	PostID int `json:"postID"`
	UserID int `json:"userID"`
}
