package contracts

import (
	"context"

	"OzonTest/internal/entity"
)

type Post interface {
	Create(ctx context.Context, input entity.Post) (*entity.Post, error)
	GetById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Post, error)
	GetAll(ctx context.Context, filter entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error)
	Delete(ctx context.Context, ID int) error
	DisableComments(ctx context.Context, ID int) (*entity.Post, error)
	EnableComments(ctx context.Context, ID int) (*entity.Post, error)
}
