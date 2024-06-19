package inmemory

import (
	"sync"
	"time"

	"OzonTest/internal/entity"
)

type comments map[int]*entity.Comment
type post struct {
	PostID    int
	UserID    int
	Content   string
	Comments  comments
	Timestamp time.Time
	IsOpen    bool

	keyGenerator int
	sync.RWMutex
}

func (p *post) insert(input entity.Post) {
	p.Lock()
	defer p.Unlock()

	*p = post{
		PostID:       p.keyGenerator,
		UserID:       input.UserID,
		Content:      input.Content,
		Comments:     make(comments),
		Timestamp:    input.Timestamp,
		IsOpen:       input.IsOpen,
		keyGenerator: 1,
	}
}

func toEntityPost(p *post) *entity.Post {
	return &entity.Post{
		PostID:    p.PostID,
		UserID:    p.UserID,
		Content:   p.Content,
		Timestamp: p.Timestamp,
		IsOpen:    p.IsOpen,
	}
}
