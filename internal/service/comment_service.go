package service

import (
	"context"
	"errors"

	"OzonTest/internal/entity"
	"OzonTest/internal/repository/contracts"
)

type CommentService struct {
	commentRepository contracts.CommentRepository
	postRepository    contracts.PostRepository
}

func NewCommentService(commentRepository contracts.CommentRepository, postRepository contracts.PostRepository) *CommentService {
	return &CommentService{commentRepository: commentRepository, postRepository: postRepository}
}

func (c *CommentService) Create(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	post, err := c.postRepository.GetById(ctx, input.PostID)
	if err != nil {
		return nil, err
	}

	if !post.IsOpen {
		return nil, errors.New("comments on the post are blocked")
	}

	return c.commentRepository.Create(ctx, input)
}

func (c *CommentService) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	return c.commentRepository.CreateRepComment(ctx, input)
}

func (c *CommentService) GetById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error) {
	return c.commentRepository.GetById(ctx, ID, pagination)
}

func (c *CommentService) GetAll(ctx context.Context, filter entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error) {
	return c.commentRepository.GetAll(ctx, filter, pagination)
}

func (c *CommentService) Delete(ctx context.Context, ID int) error {
	return c.commentRepository.Delete(ctx, ID)
}
