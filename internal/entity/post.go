package entity

import (
	"time"
)

type Post struct {
	ID        int        `json:"ID" db:"id"`
	UserID    int        `json:"userID" db:"user_id"`
	Content   string     `json:"content" db:"content"`
	Comments  []*Comment `json:"comments"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	IsOpen    bool       `json:"isOpen" db:"is_open"`
}

func (c *Post) GetTimestamp() time.Time {
	return c.Timestamp
}

type PostFilter struct {
	UserID int `json:"userID"`
}
