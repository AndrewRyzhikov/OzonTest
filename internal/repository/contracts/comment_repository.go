package contracts

import (
	"context"

	"OzonTest/internal/entity"
)

type CommentRepository interface {
	Create(ctx context.Context, input entity.Comment) (*entity.Comment, error)
	CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error)
	GetById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error)
	GetAll(ctx context.Context, filter entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error)
	Delete(ctx context.Context, ID int) error
}
