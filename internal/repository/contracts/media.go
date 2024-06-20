package contracts

import (
	"context"

	"OzonTest/internal/entity"
)

type MediaRepository interface {
	CreateComment(ctx context.Context, input entity.Comment) (*entity.Comment, error)
	CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error)
	GetCommentById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error)
	GetAllComments(ctx context.Context, filter *entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error)
	DeleteComment(ctx context.Context, ID int) error

	CreatePost(ctx context.Context, input entity.Post) (*entity.Post, error)
	GetByIdPost(ctx context.Context, ID int) (*entity.Post, error)
	GetAllPosts(ctx context.Context, filter *entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error)
	DeletePost(ctx context.Context, ID int) error
	SwitchCommentsState(ctx context.Context, ID int, state bool) (*entity.Post, error)
}
