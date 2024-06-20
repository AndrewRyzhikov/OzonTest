package entity

import (
	"time"
)

type Comment struct {
	ID        int        `json:"ID" db:"id"`
	ParentID  *int       `json:"parentID" db:"parent_id"`
	PostID    int        `json:"postID" db:"post_id"`
	UserID    int        `json:"userID" db:"user_id"`
	Content   string     `json:"content" db:"content"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Replies   []*Comment `json:"replies"`
}

func (c *Comment) GetTimestamp() time.Time {
	return c.Timestamp
}

type CommentFilter struct {
	PostID *int `json:"postID"`
	UserID *int `json:"userID"`
}
