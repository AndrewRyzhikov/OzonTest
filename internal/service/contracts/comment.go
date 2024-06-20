package contracts

import (
	"context"

	"OzonTest/internal/entity"
)

type Comment interface {
	CreateComment(ctx context.Context, input entity.Comment) (*entity.Comment, error)
	CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error)
	GetCommentById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error)
	GetAllComments(ctx context.Context, filter *entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, ID int) error
}
