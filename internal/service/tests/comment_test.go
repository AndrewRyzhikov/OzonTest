package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"OzonTest/internal/entity"
)

func TestComment_CreateComment(t *testing.T) {
	mockRepo := new(mediaRepository)
	ctx := context.TODO()
	comment := entity.Comment{ID: 1, Content: "Test Comment"}

	mockRepo.On("CreateComment", ctx, comment).Return(&comment, nil)

	createdComment, err := mockRepo.CreateComment(ctx, comment)

	require.NoError(t, err)
	require.Equal(t, &comment, createdComment)

	mockRepo.AssertExpectations(t)
}

func TestComment_CreateRepComment(t *testing.T) {
	mockRepo := new(mediaRepository)
	ctx := context.TODO()
	comment := entity.Comment{ID: 1, Content: "Reply Comment"}

	mockRepo.On("CreateRepComment", ctx, comment).Return(&comment, nil)

	createdComment, err := mockRepo.CreateRepComment(ctx, comment)

	require.NoError(t, err)
	require.Equal(t, &comment, createdComment)

	mockRepo.AssertExpectations(t)
}

func TestComment_GetCommentById(t *testing.T) {
	mockRepo := new(mediaRepository)
	ctx := context.TODO()
	commentID := 1
	limit, offser := 1, 10
	pagination := entity.Pagination{
		Limit:  &limit,
		Offset: &offser,
	}
	expectedComment := &entity.Comment{ID: commentID, Content: "Test Comment"}

	mockRepo.On("GetCommentById", ctx, commentID, pagination).Return(expectedComment, nil)

	comment, err := mockRepo.GetCommentById(ctx, commentID, pagination)

	require.NoError(t, err)
	require.Equal(t, expectedComment, comment)

	mockRepo.AssertExpectations(t)
}

func TestComment_GetAllComments(t *testing.T) {
	mockRepo := new(mediaRepository)
	ctx := context.TODO()
	filter := &entity.CommentFilter{PostID: 1}
	limit, offser := 1, 10
	pagination := entity.Pagination{
		Limit:  &limit,
		Offset: &offser,
	}
	expectedComments := []*entity.Comment{
		{ID: 1, Content: "Comment 1"},
		{ID: 2, Content: "Comment 2"},
	}

	mockRepo.On("GetAllComments", ctx, filter, pagination).Return(expectedComments, nil)

	comments, err := mockRepo.GetAllComments(ctx, filter, pagination)

	require.NoError(t, err)
	require.Equal(t, expectedComments, comments)

	mockRepo.AssertExpectations(t)
}

func TestComment_DeleteComment(t *testing.T) {
	mockRepo := new(mediaRepository)
	ctx := context.TODO()
	commentID := 1

	mockRepo.On("DeleteComment", ctx, commentID).Return(nil)

	err := mockRepo.DeleteComment(ctx, commentID)

	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
