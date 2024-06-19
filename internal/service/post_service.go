package service

import (
	"context"

	"OzonTest/internal/entity"
	"OzonTest/internal/repository/contracts"
)

type PostService struct {
	postRepository    contracts.PostRepository
	commentRepository contracts.CommentRepository
}

func NewPostService(postRepository contracts.PostRepository, commentRepository contracts.CommentRepository) *PostService {
	return &PostService{postRepository: postRepository, commentRepository: commentRepository}
}

func (p PostService) Create(ctx context.Context, input entity.Post) (*entity.Post, error) {
	return p.postRepository.Create(ctx, input)
}

func (p PostService) GetById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Post, error) {
	post, err := p.postRepository.GetById(ctx, ID)
	if err != nil {
		return &entity.Post{}, err
	}

	post.Comments, err = p.commentRepository.GetAll(ctx, entity.CommentFilter{PostID: ID}, pagination)
	if err != nil {
		return &entity.Post{}, err
	}

	return post, nil
}

func (p PostService) GetAll(ctx context.Context, filter entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error) {
	return p.postRepository.GetAll(ctx, filter, pagination)
}

func (p PostService) Delete(ctx context.Context, ID int) error {
	return p.postRepository.Delete(ctx, ID)
}

func (p PostService) DisableComments(ctx context.Context, ID int) (*entity.Post, error) {
	return p.postRepository.SwitchCommentsState(ctx, ID, false)
}

func (p PostService) EnableComments(ctx context.Context, ID int) (*entity.Post, error) {
	return p.postRepository.SwitchCommentsState(ctx, ID, true)
}
