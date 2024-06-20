package service

import (
	"context"

	"OzonTest/internal/entity"
	"OzonTest/internal/repository/contracts"
)

type PostService struct {
	repository contracts.MediaRepository
}

func NewPostService(repository contracts.MediaRepository) *PostService {
	return &PostService{repository: repository}
}

func (p *PostService) Create(ctx context.Context, input entity.Post) (*entity.Post, error) {
	return p.repository.CreatePost(ctx, input)
}

func (p *PostService) GetById(ctx context.Context, ID int, pagination entity.Pagination) (*entity.Post, error) {
	post, err := p.repository.GetByIdPost(ctx, ID)
	if err != nil {
		return &entity.Post{}, err
	}

	post.Comments, err = p.repository.GetAllComments(ctx, &entity.CommentFilter{PostID: &ID}, pagination)
	if err != nil {
		return &entity.Post{}, err
	}

	return post, nil
}

func (p *PostService) GetAll(ctx context.Context, filter *entity.PostFilter, pagination entity.Pagination) ([]*entity.Post, error) {
	return p.repository.GetAllPosts(ctx, filter, pagination)
}

func (p *PostService) Delete(ctx context.Context, ID int) error {
	return p.repository.DeletePost(ctx, ID)
}

func (p *PostService) DisableComments(ctx context.Context, ID int) (*entity.Post, error) {
	return p.repository.SwitchCommentsState(ctx, ID, false)
}

func (p *PostService) EnableComments(ctx context.Context, ID int) (*entity.Post, error) {
	return p.repository.SwitchCommentsState(ctx, ID, true)
}
