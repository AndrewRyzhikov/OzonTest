package inmemory

import (
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
}

func newPost(ID int, input entity.Post) *post {
	return &post{
		PostID:       ID,
		UserID:       input.UserID,
		Content:      input.Content,
		Comments:     make(comments),
		Timestamp:    input.Timestamp,
		IsOpen:       input.IsOpen,
		keyGenerator: 1}
}

func (p *post) insertComment(input *entity.Comment) (*entity.Comment, error) {
	if _, ok := p.Comments[input.ID]; ok {
		return &entity.Comment{}, entity.ErrAlreadyExists("comment", input.ID)
	}

	input.ID = p.keyGenerator
	p.keyGenerator++
	p.Comments[input.ID] = input

	return input, nil
}

func (p *post) insertRepComment(input *entity.Comment) (*entity.Comment, error) {
	if _, ok := p.Comments[*input.ParentID]; !ok {
		return &entity.Comment{}, entity.ErrNotFound("comment", *input.ParentID)
	}

	repComment, err := p.insertComment(input)
	if err != nil {
		return &entity.Comment{}, err
	}

	parentComment, _ := p.Comments[*input.ParentID]
	parentComment.Replies = append(parentComment.Replies, repComment)

	return repComment, nil
}

func (p *post) getAllComments() []*entity.Comment {
	comments := make([]*entity.Comment, 0, len(p.Comments))
	for _, comment := range p.Comments {
		comments = append(comments, comment)
	}

	return comments
}

func (p *post) findComment(id int) (*entity.Comment, bool) {
	comment, ok := p.Comments[id]
	return comment, ok
}

func (p *post) deleteComment(id int) error {
	comment, ok := p.Comments[id]
	if !ok {
		return entity.ErrNotFound("comment", id)
	}

	for _, repComment := range comment.Replies {
		p.deleteComment(repComment.ID)
	}

	delete(p.Comments, id)

	return nil
}

func toEntityPost(p *post) *entity.Post {
	comments := make([]*entity.Comment, 0, len(p.Comments))

	for _, comment := range p.Comments {
		comments = append(comments, comment)
	}

	return &entity.Post{
		ID:        p.PostID,
		UserID:    p.UserID,
		Content:   p.Content,
		Timestamp: p.Timestamp,
		Comments:  comments,
		IsOpen:    p.IsOpen,
	}
}
