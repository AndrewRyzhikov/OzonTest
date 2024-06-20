package service

import (
	"context"
	"errors"

	"OzonTest/internal/entity"
	"OzonTest/internal/repository/contracts"
)

type CommentService struct {
	repository contracts.MediaRepository
}

func NewCommentService(repository contracts.MediaRepository) *CommentService {
	return &CommentService{repository: repository}
}

func (c *CommentService) CreateComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	post, err := c.repository.GetByIdPost(ctx, input.PostID)
	if err != nil {
		return &entity.Comment{}, err
	}

	if !post.IsOpen {
		return &entity.Comment{}, errors.New("comments on the post are blocked")
	}

	return c.repository.CreateComment(ctx, input)
}

func (c *CommentService) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	post, err := c.repository.GetByIdPost(ctx, input.PostID)
	if err != nil {
		return &entity.Comment{}, err
	}

	if !post.IsOpen {
		return &entity.Comment{}, errors.New("comments on the post are blocked")
	}

	return c.repository.CreateRepComment(ctx, input)
}

func (c *CommentService) GetCommentById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error) {
	return c.repository.GetCommentById(ctx, ID, pagination)
}

func (c *CommentService) GetAllComments(ctx context.Context, filter *entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error) {
	return c.repository.GetAllComments(ctx, filter, pagination)
}

func (c *CommentService) DeleteComment(ctx context.Context, ID int) error {
	return c.repository.DeleteComment(ctx, ID)
}
