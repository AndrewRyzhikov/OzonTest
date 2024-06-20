package tests

import (
	"context"

	"github.com/stretchr/testify/mock"

	"OzonTest/internal/entity"
)

type mediaRepository struct {
	mock.Mock
}

func (m *mediaRepository) CreateComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *mediaRepository) CreateRepComment(ctx context.Context, input entity.Comment) (*entity.Comment, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *mediaRepository) GetCommentById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Comment, error) {
	args := m.Called(ctx, ID, pagination)
	return args.Get(0).(*entity.Comment), args.Error(1)
}

func (m *mediaRepository) GetAllComments(ctx context.Context, filter *entity.CommentFilter, pagination entity.Pagination) ([]*entity.Comment, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).([]*entity.Comment), args.Error(1)
}

func (m *mediaRepository) DeleteComment(ctx context.Context, ID int) error {
	args := m.Called(ctx, ID)
	return args.Error(0)
}

func (m *mediaRepository) CreatePost(ctx context.Context, input entity.Post) (*entity.Post, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *mediaRepository) GetByIdPost(ctx context.Context, ID int) (*entity.Post, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(*entity.Post), args.Error(1)
}

func (m *mediaRepository) GetAllPosts(ctx context.Context, filter *entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error) {
	args := m.Called(ctx, filter, pagination)
	return args.Get(0).([]*entity.Post), args.Error(1)
}

func (m *mediaRepository) DeletePost(ctx context.Context, ID int) error {
	args := m.Called(ctx, ID)
	return args.Error(0)
}

func (m *mediaRepository) SwitchCommentsState(ctx context.Context, ID int, state bool) (*entity.Post, error) {
	args := m.Called(ctx, ID, state)
	return args.Get(0).(*entity.Post), args.Error(1)
}
