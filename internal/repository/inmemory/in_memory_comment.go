package inmemory

import (
	"context"
	"errors"
	"sync"

	"OzonTest/internal/entity"
)

type comments map[int]*entity.Comment

type CommentInMemory struct {
	comments     comments
	keyGenerator int
	sync.RWMutex
}

func NewCommentInMemory() *CommentInMemory {
	return &CommentInMemory{comments: make(comments), keyGenerator: 1}
}

func (c *CommentInMemory) Create(_ context.Context, input entity.Comment) (*entity.Comment, error) {
	c.Lock()
	defer c.Unlock()

	input.ID = c.keyGenerator
	c.comments[c.keyGenerator] = &input

	c.keyGenerator++

	return &input, nil
}

func (c *CommentInMemory) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	c.Lock()

	if _, ok := c.comments[input.ParentID]; !ok {
		return &entity.Comment{}, errors.New("parent ID not found")
	}

	c.Unlock()

	return c.Create(ctx, input)
}

func (c *CommentInMemory) GetById(_ context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error) {
	c.RLock()
}

func (c *CommentInMemory) GetAll(_ context.Context, filter entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CommentInMemory) Delete(ctx context.Context, ID int) error {
	//TODO implement me
	panic("implement me")
}
