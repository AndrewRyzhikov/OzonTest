package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"OzonTest/internal/entity"
	"OzonTest/internal/service"
)

func TestPostService_Create(t *testing.T) {
	mockRepo := new(mediaRepository)
	postService := service.NewPostService(mockRepo)
	ctx := context.TODO()
	post := entity.Post{ID: 1, Content: "Test Post"}

	mockRepo.On("CreatePost", ctx, post).Return(&post, nil)

	createdPost, err := postService.Create(ctx, post)

	require.NoError(t, err)
	require.Equal(t, &post, createdPost)

	mockRepo.AssertExpectations(t)
}

func TestPostService_GetById(t *testing.T) {
	mockRepo := new(mediaRepository)
	postService := service.NewPostService(mockRepo)
	ctx := context.TODO()
	postID := 1

	limit, offser := 1, 10
	pagination := entity.Pagination{
		Limit:  &limit,
		Offset: &offser,
	}
	expectedPost := &entity.Post{ID: postID, Content: "Test Post"}
	expectedComments := []*entity.Comment{
		{ID: 1, PostID: postID, Content: "Comment 1"},
		{ID: 2, PostID: postID, Content: "Comment 2"},
	}

	mockRepo.On("GetByIdPost", ctx, postID).Return(expectedPost, nil)

	mockRepo.On("GetAllComments", ctx, mock.AnythingOfType("*entity.CommentFilter"), pagination).Return(expectedComments, nil)

	post, err := postService.GetById(ctx, postID, pagination)

	require.NoError(t, err)
	require.Equal(t, expectedPost, post)
	require.Equal(t, expectedComments, post.Comments)

	mockRepo.AssertExpectations(t)
}

func TestPostService_GetAll(t *testing.T) {
	mockRepo := new(mediaRepository)
	postService := service.NewPostService(mockRepo)
	ctx := context.TODO()
	filter := &entity.PostFilter{}

	limit, offser := 1, 10
	pagination := entity.Pagination{
		Limit:  &limit,
		Offset: &offser,
	}

	expectedPosts := []*entity.Post{
		{ID: 1, Content: "Test Post 1"},
		{ID: 2, Content: "Test Post 2"},
	}

	mockRepo.On("GetAllPosts", ctx, filter, pagination).Return(expectedPosts, nil)

	posts, err := postService.GetAll(ctx, filter, pagination)

	require.NoError(t, err)
	require.Equal(t, expectedPosts, posts)

	mockRepo.AssertExpectations(t)
}

func TestPostService_Delete(t *testing.T) {
	mockRepo := new(mediaRepository)
	postService := service.NewPostService(mockRepo)
	ctx := context.TODO()
	postID := 1

	mockRepo.On("DeletePost", ctx, postID).Return(nil)

	err := postService.Delete(ctx, postID)

	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestPostService_DisableComments(t *testing.T) {
	mockRepo := new(mediaRepository)
	postService := service.NewPostService(mockRepo)
	ctx := context.TODO()
	postID := 1
	expectedPost := &entity.Post{ID: postID, Content: "Test Post", IsOpen: false}

	mockRepo.On("SwitchCommentsState", ctx, postID, false).Return(expectedPost, nil)

	post, err := postService.DisableComments(ctx, postID)

	require.NoError(t, err)
	require.Equal(t, expectedPost, post)

	mockRepo.AssertExpectations(t)
}

func TestPostService_EnableComments(t *testing.T) {
	mockRepo := new(mediaRepository)
	postService := service.NewPostService(mockRepo)
	ctx := context.TODO()
	postID := 1
	expectedPost := &entity.Post{ID: postID, Content: "Test Post", IsOpen: true}

	mockRepo.On("SwitchCommentsState", ctx, postID, true).Return(expectedPost, nil)

	post, err := postService.EnableComments(ctx, postID)

	require.NoError(t, err)
	require.Equal(t, expectedPost, post)

	mockRepo.AssertExpectations(t)
}
